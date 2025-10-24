package main

import (
	"fmt"
	"log"
	"os"

	vh "github.com/Tryanks/virtualhere-go"
)

func main() {
	// Check for required arguments
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run use_device.go <device_address> [password]")
		fmt.Println("Example: go run use_device.go raspberrypi.114")
		os.Exit(1)
	}

	deviceAddress := os.Args[1]
	password := ""
	if len(os.Args) > 2 {
		password = os.Args[2]
	}

	// Create VirtualHere pipe client (connects to running service)
	client, err := vh.NewPipeClient()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	fmt.Println("VirtualHere - Use Device Example")
	fmt.Println("=================================\n")

	// Use the device
	fmt.Printf("Connecting to device: %s\n", deviceAddress)
	err = client.Use(deviceAddress, password)
	if err != nil {
		log.Fatalf("Failed to use device: %v", err)
	}

	fmt.Println("✓ Device connected successfully!")
	fmt.Println("\nPress Enter to disconnect the device...")
	fmt.Scanln()

	// Stop using the device
	fmt.Println("Disconnecting device...")
	err = client.StopUsing(deviceAddress)
	if err != nil {
		log.Fatalf("Failed to stop using device: %v", err)
	}

	fmt.Println("✓ Device disconnected successfully!")
}
