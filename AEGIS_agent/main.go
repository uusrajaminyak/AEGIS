package main

import (
"context"
	"fmt"
	"log"
	"syscall"
	"unsafe"
	"time"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "github.com/uusrajaminyak/aegis-backend/api/proto"
)

var hqClient pb.AegisSentinelClient

func cStringToGo(ptr uintptr) string {
		if ptr == 0 {
				return ""
		}
		var bytes []byte
		basePtr := *(*unsafe.Pointer)(unsafe.Pointer(&ptr))
		for i := 0; ; i++ {
				b := *(*byte)(unsafe.Add(basePtr, i))
				if b == 0 {
						break
				}
				bytes = append(bytes, b)
		}
		return string(bytes)
}

func onAlertReceived(severity uintptr, messagePtr uintptr) uintptr {
		message := cStringToGo(messagePtr)
		fmt.Printf("Alert received - Severity: %d, Message: %s\n", severity, message)
		fmt.Printf("Sending alert to HQ...\n")

		if hqClient != nil {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				req := &pb.AlertRequest{
						AgentId: "agent-123",
						Description: message,
						EventType: "SensorAlert",
						Severity: fmt.Sprintf("%d", severity),
				}
				_, err := hqClient.SendAlert(ctx, req)
				if err != nil {
						log.Printf("Failed to send alert to HQ: %v", err)
				} else {
						log.Printf("Alert sent to HQ successfully")
				}
		}
		return 0
}

func main() {
		fmt.Println("Loading sensor module...")
		dllPath := "core/aegis_core.dll"
		conn, err := grpc.Dial("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
				log.Fatalf("Failed to connect to HQ: %v", err)
		}
		defer conn.Close()
		hqClient = pb.NewAegisSentinelClient(conn)
		fmt.Println("Connected to HQ successfully.")
		aegisCore, err := syscall.LoadDLL(dllPath)
		if err != nil {
				log.Fatalf("Failed to load DLL: %v", err)
		}
		defer aegisCore.Release()

		setCallbackProc, err := aegisCore.FindProc("SetAlertCallback")
		if err != nil {
				log.Fatalf("Failed to find SetAlertCallback procedure: %v", err)
		}

		callback := syscall.NewCallback(onAlertReceived)
		setCallbackProc.Call(callback)

		fmt.Println("DLL loaded successfully.")

		initSensor, err := aegisCore.FindProc("InitSensor")
		if err != nil {
				log.Fatalf("Failed to find InitSensor procedure: %v", err)
		}
		initSensor.Call()
		fmt.Println("Sensor initialized successfully.")
}