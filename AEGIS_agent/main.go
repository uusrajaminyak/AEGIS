package main

import (
		"fmt"
		"log"
		"syscall"
		"unsafe"
)

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

func onAlertReceived(severity int, messagePtr uintptr) uintptr {
		message := cStringToGo(messagePtr)
		fmt.Printf("Alert received - Severity: %d, Message: %s\n", severity, message)
		return 0
}

func main() {
		fmt.Println("Loading sensor module...")
		dllPath := "core/aegis_core.dll"
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