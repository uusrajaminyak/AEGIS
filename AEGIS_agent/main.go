package main

import (
		"fmt"
		"log"
		"syscall"
)

func main() {
		fmt.Println("Loading sensor module...")
		dllPath := "core/aegis_core.dll"
		aegisCore, err := syscall.LoadDLL(dllPath)
		if err != nil {
				log.Fatalf("Failed to load DLL: %v", err)
		}
		defer aegisCore.Release()

		fmt.Println("DLL loaded successfully.")

		initSensor, err := aegisCore.FindProc("InitSensor")
		if err != nil {
				log.Fatalf("Failed to find InitSensor procedure: %v", err)
		}
		initSensor.Call()
		fmt.Println("Sensor initialized successfully.")
}