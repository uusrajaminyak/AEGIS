package main

import (
"context"
	"fmt"
	"log"
	"os/exec"
	"syscall"
	"unsafe"
	"strings"
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
						EventType: "CreateProcess_Hook",
						Severity: fmt.Sprintf("%d", severity),
				}
				res, err := hqClient.SendAlert(ctx, req)
				if err != nil {
						log.Printf("Failed to send alert to HQ: %v", err)
				} else {
						log.Printf("Alert sent to HQ successfully")
						log.Printf("HQ Response - Alert ID: %s, Action: %s, Target: %s", res.AlertId, res.Action, res.Target)
						if res.Action == "KILL" && res.Target != "" {
								go func(targetToKill string) {
										time.Sleep(1 * time.Second)
										killProcessByPattern(targetToKill)
								}(res.Target)
						}
				}
		}
		return 0
}

func killProcessByPattern(pattern string) {
		fmt.Printf("Hunting process matching: '%s'", pattern)
		safeRegex := strings.ReplaceAll(pattern, " ", ".*")
		psCmd := fmt.Sprintf(`Get-WmiObject Win32_Process | Where-Object { $_.ProcessId -ne $PID -and $_.CommandLine -match '%s' } | ForEach-Object { Stop-Process -Id $_.ProcessId -Force }`, safeRegex)
		cmd := exec.Command("powershell", "-Command", psCmd)
		err := cmd.Run()
		if err != nil {
				log.Printf("Failed to kill process pattern '%s': %v", pattern, err)
		} else {
				log.Printf("Process pattern '%s' killed successfully", pattern)
		}
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

		testNetProc, err := aegisCore.FindProc("TestNetworkHook")
		if err != nil {
				log.Fatalf("Failed to find TestNetworkHook procedure: %v", err)
		}

		fmt.Println("Testing network hook...")
		fmt.Println("Simulating fileless process creation...")

		go func() {
				time.Sleep(2 * time.Second)
				cmd := exec.Command("powershell.exe", "-WindowStyle", "Hidden", "-ExecutionPolicy", "Bypass", "-Command", "Start-Sleep -Seconds 15")
				err := cmd.Start()
				if err != nil {
						log.Printf("Failed to simulate fileless malware: %v", err)
				} else {
						log.Printf("Fileless malware simulation started with PID: %d", cmd.Process.Pid)
				}

				time.Sleep(2 * time.Second)
				testNetProc.Call()
		}()
		select {}
}