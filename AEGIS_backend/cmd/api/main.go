package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	pb "github.com/uusrajaminyak/aegis-backend/api/proto"
	"github.com/uusrajaminyak/aegis-backend/config"
	"github.com/uusrajaminyak/aegis-backend/internal/adapter"
	grpc_handler "github.com/uusrajaminyak/aegis-backend/internal/handler/grpc"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"os"
	"sync"
)

type RateLimiterManager struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
}

type CommandPayload struct {
	AgentID       string `json:"agent_id"`
	Command       string `json:"command"`
	TargetProcess string `json:"target_process"`
}

func NewRateLimiterManager(r rate.Limit, b int) *RateLimiterManager {
	return &RateLimiterManager{
		limiters: make(map[string]*rate.Limiter),
		rate:     r,
		burst:    b,
	}
}

func (rlm *RateLimiterManager) GetLimiter(clientID string) *rate.Limiter {
	rlm.mu.RLock()
	limiter, exists := rlm.limiters[clientID]
	rlm.mu.RUnlock()
	if !exists {
		rlm.mu.Lock()
		limiter, exists = rlm.limiters[clientID]
		if !exists {
			limiter = rate.NewLimiter(rlm.rate, rlm.burst)
			rlm.limiters[clientID] = limiter
		}
		rlm.mu.Unlock()
	}
	return limiter
}

func (rlm *RateLimiterManager) UnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var AgentID string
		if alertReq, ok := req.(*pb.AlertRequest); ok {
			AgentID = alertReq.AgentId
		} else {
			AgentID = "unknown"
		}
		if AgentID != "unknown" {
			limiter := rlm.GetLimiter(AgentID)
			if !limiter.Allow() {
				return nil, status.Errorf(codes.ResourceExhausted, "rate limit exceeded for agent %s", AgentID)
			}
		}
		return handler(ctx, req)
	}
}

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("[!] Failed to load config: %v", err)
	}
	adapter.ConnectPostgres(cfg)

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("[!] Failed to connect to NATS server: %v", err)
	}
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Fatalf("[!] Failed to get JetStream context: %v", err)
	}

	_, err = js.AddStream(&nats.StreamConfig{
		Name:     "Alerts_Stream",
		Subjects: []string{"Alerts.*"},
	})
	if err != nil {
		log.Printf("[!] Failed to add stream (might already exist): %v", err)
	} else {
		log.Printf("[*] Stream 'Alerts_Stream' created successfully")
	}

	go func() {
		log.Println("[*] Starting NATS worker to process alerts...")
		_, err := js.Subscribe("Alerts.new", func(msg *nats.Msg) {
			var payload map[string]interface{}
			if err := json.Unmarshal(msg.Data, &payload); err != nil {
				log.Printf("[!] Failed to unmarshal alert message: %v", err)
				return
			}

			agentID := payload["agent_id"].(string)
			eventType := payload["event_type"].(string)
			severity := payload["severity"].(string)
			description := payload["description"].(string)
			action := payload["action"].(string)
			targetProcess := payload["target_process"].(string)

			query := `INSERT INTO alerts (agent_id, event_type, severity, description, action, target_process) VALUES ($1, $2, $3, $4, $5, $6)`
			result := adapter.DB.Exec(query, agentID, eventType, severity, description, action, targetProcess)
			if result.Error != nil {
				log.Printf("[!] Failed to insert alert into database: %v", result.Error)
			} else {
				log.Printf("[+] Alert from agent %s stored in database successfully", agentID)
			}
		}, nats.DeliverNew())
		if err != nil {
			log.Printf("[!] Failed to subscribe to NATS subject: %v", err)
		}
	}()

	go func() {
		grpcPort := ":9090"
		lis, err := net.Listen("tcp", grpcPort)
		if err != nil {
			log.Fatalf("[!] Failed to listen on port %s: %v", grpcPort, err)
		}
		caCert, err := os.ReadFile("cert/ca.crt")
		if err != nil {
			log.Fatalf("[!] Failed to read CA certificate: %v", err)
		}
		certPool := x509.NewCertPool()
		if !certPool.AppendCertsFromPEM(caCert) {
			log.Fatalf("[!] Failed to append CA certificate to pool")
		}
		serverCert, err := tls.LoadX509KeyPair("cert/server.crt", "cert/server.key")
		if err != nil {
			log.Fatalf("[!] Failed to load server certificate and key: %v", err)
		}
		creds := credentials.NewTLS(&tls.Config{
			Certificates: []tls.Certificate{serverCert},
			ClientCAs:    certPool,
			ClientAuth:   tls.RequireAndVerifyClientCert,
			MinVersion:   tls.VersionTLS13,
		})

		limitManager := NewRateLimiterManager(rate.Limit(5), 10)

		grpcServer := grpc.NewServer(grpc.Creds(creds), grpc.UnaryInterceptor(limitManager.UnaryInterceptor()))

		sentinelHandler := &grpc_handler.SentinelServer{
			DB: adapter.DB,
			JS: js,
		}
		pb.RegisterAegisSentinelServer(grpcServer, sentinelHandler)
		fmt.Printf("[*] gRPC server listening on %s\n", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("[!] Failed to serve gRPC server: %v", err)
		}
	}()
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/api/alerts", func(c *gin.Context) {
		var alerts []map[string]interface{}
		result := adapter.DB.Table("alerts").Order("id Desc").Find(&alerts)
		if result.Error != nil {
			c.JSON(500, gin.H{"error": "Failed to fetch alerts"})
			return
		}

		c.JSON(200, gin.H{
			"status": "success",
			"data":   alerts,
		})
	})

	r.POST("/api/command", func(c *gin.Context) {
		var payload CommandPayload

		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request payload"})
			return
		}

		fmt.Printf("[*] Received command request: AgentID=%s, Command=%s, TargetProcess=%s\n", payload.AgentID, payload.Command, payload.TargetProcess)

		commandBytes, err := json.Marshal(payload)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to marshal command payload"})
			return
		}

		natsTopic := "agent." + payload.AgentID + ".command"
		err = nc.Publish(natsTopic, commandBytes)

		if err != nil {
			log.Printf("[!] Failed to publish command to NATS: %v", err)
			c.JSON(500, gin.H{"error": "Failed to send command to agent"})
			return
		}

		c.JSON(200, gin.H{
			"status": "success",
			"message":   "Command: " + payload.Command + " sent to agent " + payload.AgentID,
		})
	})

	httpPort := cfg.ServerPort
	if httpPort == "" {
		httpPort = ":8080"
	}

	fmt.Printf("[*] HTTP server listening on %s\n", httpPort)
	if err := r.Run(":" + httpPort); err != nil {
		log.Fatalf("[!] Failed to run HTTP server: %v", err)
	}
}