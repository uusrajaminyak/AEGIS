package grpc

import (
	"context"
	"github.com/google/uuid"
	pb "github.com/uusrajaminyak/aegis-backend/api/proto"
	"gorm.io/gorm"
	"log"
	"strings"
	"time"
	"encoding/json"
	"github.com/nats-io/nats.go"
)

type SentinelServer struct {
	pb.UnimplementedAegisSentinelServer
	DB *gorm.DB
	JS nats.JetStreamContext
}

type DetectionRule struct {
	ID          uint   `gorm:"primaryKey"`
	ProcessName string `gorm:"unique;not null"`
	Action      string `gorm:"not null"`
	IsActive    bool   `gorm:"default:true"`
}

func (s *SentinelServer) Connect(ctx context.Context, req *pb.ConnectRequest) (*pb.ConnectResponse, error) {
	log.Printf("[+] Received Connect request from New Agent: %s (%s)", req.Hostname, req.IpAddress)
	newID := uuid.New().String()
	log.Printf("[*] Assigned Agent ID: %s", newID)
	return &pb.ConnectResponse{
		AgentId:   newID,
		Status:    "Connected",
		AuthToken: "auth-" + newID,
	}, nil
}

func (s *SentinelServer) SendHeartbeat(ctx context.Context, req *pb.HeartbeatRequest) (*pb.HeartbeatResponse, error) {
	return &pb.HeartbeatResponse{
		Status: "Ok",
	}, nil
}

func (s *SentinelServer) SendAlert(ctx context.Context, req *pb.AlertRequest) (*pb.AlertResponse, error) {
	log.Printf("[+] Received Alert from %s: ", req.AgentId)
	log.Printf("[*] Event type: %s", req.EventType)
	log.Printf("[*] Severity: %s", req.Severity)
	log.Printf("[*] Alert Details: %s", req.Description)

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
				log.Printf("[+] Detection rule matched: %s", badProcess)
				log.Printf("[*] Action: %s", action)
				break
			}
		}
	}

	if s.JS != nil {
		payload := map[string]interface{}{
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"agent_id": req.AgentId,
			"event_type": req.EventType,
			"severity": req.Severity,
			"description": req.Description,
			"action": action,
			"target_process": targetProcess,
		}
		jsonData, err := json.Marshal(payload)
		if err != nil {
			log.Printf("[!] Failed to encode alert to JSON: %v", err)
		} else {
			_, err = s.JS.Publish("Alerts.new", jsonData)
			if err != nil {
				log.Printf("[!] Failed to publish alert to NATS: %v", err)
			} else {
				log.Printf("[+] Alert published to NATS Queue")
			}
		}
	} else {
		log.Printf("[!] NATS JetStream context not available, skipping alert publication")
	}

	return &pb.AlertResponse{
		AlertId: uuid.New().String(),
		Action:  action,
		Target:  targetProcess,
	}, nil
}

func (s *SentinelServer) CommandStream(req *pb.CommandRequest, stream pb.AegisSentinel_CommandStreamServer) error {
	log.Printf("[+] Agent %s connected for command stream", req.AgentId)
	for {
		var activeRules []DetectionRule
		if s.DB != nil {
			s.DB.Where("is_active = ?", true).Find(&activeRules)
		}

		var ruleNames []string
		for _, rule := range activeRules {
			ruleNames = append(ruleNames, strings.ToLower(rule.ProcessName))
		}
		combinedRules := strings.Join(ruleNames, ",")
		res := &pb.CommandMessage{
			Type:    "SYNC_RULES",
			Payload: combinedRules,
		}

		if err := stream.Send(res); err != nil {
			log.Printf("[!] Error sending command to agent %s: %v", req.AgentId, err)
			return err
		}
		time.Sleep(30 * time.Second)
	}
}