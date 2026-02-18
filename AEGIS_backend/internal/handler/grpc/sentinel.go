package grpc

import (
		"context"
		"gorm.io/gorm"
		"log"
		"time"
		"strings"
		"github.com/google/uuid"
		pb "github.com/uusrajaminyak/aegis-backend/api/proto"
)

type SentinelServer struct {
		pb.UnimplementedAegisSentinelServer
		DB *gorm.DB
}

type DetectionRule struct {
		ID uint `gorm:"primaryKey"`
		ProcessName string `gorm:"unique;not null"`
		Action string `gorm:"not null"`
		IsActive bool `gorm:"default:true"`
}

func (s *SentinelServer) Connect(ctx context.Context, req *pb.ConnectRequest) (*pb.ConnectResponse, error) {
		log.Printf("Received Connect request from New Agent: %s (%s)", req.Hostname, req.IpAddress)
		newID := uuid.New().String()
		log.Printf("Assigned Agent ID: %s", newID)
		return &pb.ConnectResponse{
				AgentId: newID,
				Status: "Connected",
				AuthToken: "auth-" + newID,
		}, nil
}

func (s *SentinelServer) SendHeartbeat(ctx context.Context, req *pb.HeartbeatRequest) (*pb.HeartbeatResponse, error) {
		return &pb.HeartbeatResponse{
				Status: "Ok",
		}, nil
}

func (s *SentinelServer) SendAlert(ctx context.Context, req *pb.AlertRequest) (*pb.AlertResponse, error) {
		log.Printf("Received Alert from %s: ", req.AgentId)
		log.Printf("Event type: %s", req.EventType)
		log.Printf("Severity: %s", req.Severity)
		log.Printf("Alert Details: %s", req.Description)

		action := "Log_and_Investigate"
		targetProcess := ""
		lowerMsg := strings.ToLower(req.Description)

		if req.EventType == "CreateProcess_Hook" && !strings.Contains(lowerMsg, "taskkill") {
				var activeRules []DetectionRule
				if s.DB != nil {
						s.DB.AutoMigrate(&DetectionRule{})
						s.DB.Where("is_active = ?", true).Find(&activeRules)
				}
				for _, rule := range activeRules {
						badProcess := strings.ToLower(rule.ProcessName)
						if strings.Contains(lowerMsg, badProcess) {
								action = rule.Action
								targetProcess = rule.ProcessName
								log.Printf("Detection rule matched: %s, Action: %s", badProcess, action)
								break
						}
				}
		}

		if s.DB != nil {
				query := `INSERT INTO alerts (agent_id, event_type, severity, description) VALUES (?, ?, ?, ?)`
				result := s.DB.Exec(query, req.AgentId, req.EventType, req.Severity, req.Description)
				if result.Error != nil {
						log.Printf("Failed to store alert in database: %v", result.Error)
				} else {
						log.Printf("Alert stored in database successfully")
				}
		} else {
				log.Printf("Database connection not available, skipping alert storage")
		}

		return &pb.AlertResponse{
				AlertId: uuid.New().String(),
				Action: action,
				Target: targetProcess,
		}, nil
}

func (s *SentinelServer) CommandStream(req *pb.CommandRequest, stream pb.AegisSentinel_CommandStreamServer) error {
		for {
				time.Sleep(10 * time.Second)
		}
} 