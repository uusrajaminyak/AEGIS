package main

import (
"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
	"unsafe"
	"strings"
	"time"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"github.com/shirou/gopsutil/v3/process"
	"golang.org/x/sys/windows"
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
						log.Printf("Failed to send alert to HQ: %v\n", err)
				} else {
						log.Printf("Alert sent to HQ successfully\n")
						log.Printf("HQ Response - Alert ID: %s, Action: %s, Target: %s\n", res.AlertId, res.Action, res.Target)
						if res.Action == "KILL" && res.Target != "" {
								go func(targetToKill string) {
										time.Sleep(1 * time.Second)
										killProcessNative(targetToKill)
								}(res.Target)
						}
				}
		}
		return 0
}

func killProcessNative(pattern string) {
		fmt.Printf("scanning for processes matching pattern: %s\n", pattern)
		processes, err := process.Processes()
		if err != nil {
				log.Printf("Failed to list processes: %v\n", err)
				return
		}

		targetPID := int32(0)
		safePattern := strings.ToLower(pattern)

		for _, p := range processes {
				if int32(os.Getpid()) == p.Pid {
						continue
				}
				cmdline, err := p.Cmdline()
				if err == nil && strings.Contains(strings.ToLower(cmdline), safePattern) {
						targetPID = p.Pid
						log.Printf("Found matching process - PID: %d, Cmdline: %s\n", p.Pid, cmdline)
						break
				}
		}

		if targetPID == 0 {
				log.Printf("No process found matching pattern: %s\n", pattern)
				return
		}

		const PROCESS_TERMINATE = 0x0001
		handle, err := windows.OpenProcess(PROCESS_TERMINATE, false, uint32(targetPID))
		if err != nil {
				log.Printf("Failed to open process with PID %d: %v\n", targetPID, err)
				return
		}

		defer windows.CloseHandle(handle)
		err = windows.TerminateProcess(handle, 1)
		if err != nil {
				log.Printf("Failed to terminate process with PID %d: %v\n", targetPID, err)
		} else {
				log.Printf("Successfully terminated process with PID %d\n", targetPID)
		}
}

func main() {
		fmt.Println("Loading sensor module...")
		dllPath := "core/aegis_core.dll"
		conn, err := grpc.Dial("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
				log.Fatalf("Failed to connect to HQ: %v\n", err)
		}
		defer conn.Close()
		hqClient = pb.NewAegisSentinelClient(conn)
		fmt.Println("Connected to HQ successfully.")
		aegisCore, err := syscall.LoadDLL(dllPath)
		if err != nil {
				log.Fatalf("Failed to load DLL: %v\n", err)
		}
		defer aegisCore.Release()

		setCallbackProc, err := aegisCore.FindProc("SetAlertCallback")
		if err != nil {
				log.Fatalf("Failed to find SetAlertCallback procedure: %v\n", err)
		}

		callback := syscall.NewCallback(onAlertReceived)
		setCallbackProc.Call(callback)

		fmt.Println("DLL loaded successfully.")

		initSensor, err := aegisCore.FindProc("InitSensor")
		if err != nil {
				log.Fatalf("Failed to find InitSensor procedure: %v\n", err)
		}
		initSensor.Call()
		fmt.Println("Sensor initialized successfully.")

		testNetProc, err := aegisCore.FindProc("TestNetworkHook")
		if err != nil {
				log.Fatalf("Failed to find TestNetworkHook procedure: %v\n", err)
		}

		fmt.Println("Testing network hook...")
		fmt.Println("Simulating fileless process creation...")

		go func() {
				time.Sleep(2 * time.Second)
				cmd := exec.Command("powershell.exe", "-WindowStyle", "Hidden", "-ExecutionPolicy", "Bypass", "-Command", "Start-Sleep -Seconds 15")
				err := cmd.Start()
				if err != nil {
						log.Printf("Failed to simulate fileless malware: %v\n", err)
				} else {
						log.Printf("Fileless malware simulation started with PID: %d\n", cmd.Process.Pid)
				}

				time.Sleep(2 * time.Second)
				testNetProc.Call()
		}()
		select {}
}