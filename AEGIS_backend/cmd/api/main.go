package main

import (
	"fmt"
	"log"
	"net"

	"github.com/gin-gonic/gin"
	pb "github.com/uusrajaminyak/aegis-backend/api/proto"
	"github.com/uusrajaminyak/aegis-backend/config"
	"github.com/uusrajaminyak/aegis-backend/internal/adapter"
	grpc_handler "github.com/uusrajaminyak/aegis-backend/internal/handler/grpc"
	"google.golang.org/grpc"
)

func main() {
		cfg, err := config.LoadConfig(".")
		if err != nil {
				log.Fatalf("Failed to load config: %v", err)
		}
		adapter.ConnectPostgres(cfg)
		go func() {
				grpcPort := ":9090"
				lis, err := net.Listen("tcp", grpcPort)
				if err != nil {
						log.Fatalf("Failed to listen on port %s: %v", grpcPort, err)
				}
				grpcServer := grpc.NewServer()
				sentinelHandler := &grpc_handler.SentinelServer{}
				pb.RegisterAegisSentinelServer(grpcServer, sentinelHandler)
				fmt.Printf("gRPC server listening on %s\n", grpcPort)
				if err := grpcServer.Serve(lis); err != nil {
						log.Fatalf("Failed to serve gRPC server: %v", err)
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

		fmt.Printf("HTTP server listening on %s\n", httpPort)
		if err := r.Run(":" + httpPort); err != nil {
				log.Fatalf("Failed to run HTTP server: %v", err)
		}
}