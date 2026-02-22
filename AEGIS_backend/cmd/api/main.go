package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/gin-gonic/gin"
	pb "github.com/uusrajaminyak/aegis-backend/api/proto"
	"github.com/uusrajaminyak/aegis-backend/config"
	"github.com/uusrajaminyak/aegis-backend/internal/adapter"
	grpc_handler "github.com/uusrajaminyak/aegis-backend/internal/handler/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
	"os"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("[!] Failed to load config: %v", err)
	}
	adapter.ConnectPostgres(cfg)
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
		grpcServer := grpc.NewServer(grpc.Creds(creds))

		sentinelHandler := &grpc_handler.SentinelServer{DB: adapter.DB}
		pb.RegisterAegisSentinelServer(grpcServer, sentinelHandler)
		fmt.Printf("[*] gRPC server listening on %s\n", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("[!] Failed to serve gRPC server: %v", err)
		}
	}()
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
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