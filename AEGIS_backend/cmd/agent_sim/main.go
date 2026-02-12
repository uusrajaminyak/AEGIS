package main

import (
		"fmt"
		"log"
		"context"
		"time"

		"google.golang.org/grpc"
		"google.golang.org/grpc/credentials/insecure"
		pb "github.com/uusrajaminyak/aegis-backend/api/proto"
)

func main() {
		serverAddr := "localhost:9090"
		conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
				log.Fatalf("Failed to connect to gRPC server: %v", err)
		}
		defer conn.Close()

		client := pb.NewAegisSentinelClient(conn)
		fmt.Println("Connecting to Sentinel gRPC server...")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		response, err := client.Connect(ctx, &pb.ConnectRequest{
				Hostname:  "agent-sim-01",
				IpAddress: "192.168.1.105",
				OsVersion: "Windows 11",
				AgentVersion: "v0.1",
				PublicKey: "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQE",
		})

		if err != nil {
				log.Fatalf("Failed to connect to Sentinel: %v", err)
		}

		fmt.Printf("Connected! Assigned Agent ID: %s, Auth Token: %s\n", response.AgentId, response.AuthToken)
		fmt.Printf("Status: %s\n", response.Status)
}