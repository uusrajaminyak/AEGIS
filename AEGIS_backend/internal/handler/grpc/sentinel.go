package grpc

import (
		"context"
		"log"
		"time"
		"github.com/google/uuid"
		pb "github.com/uusrajaminyak/aegis-backend/api/proto"
)

type SentinelServer struct {
		pb.UnimplementedAegisSentinelServer
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
		log.Printf("Received Alert from Agent %s: %s (%s)", req.AgentId, req.EventType, req.Severity)
		log.Printf("Alert Details: %s", req.Description)

		return &pb.AlertResponse{
				AlertId: uuid.New()	.String(),
				Action: "Log_only",
		}, nil
}

func (s *SentinelServer) CommandStream(req *pb.CommandRequest, stream pb.AegisSentinel_CommandStreamServer) error {
		for {
				time.Sleep(10 * time.Second)
		}
} 