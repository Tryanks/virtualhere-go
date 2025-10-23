package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	vh "github.com/Tryanks/virtualhere-go"
)

func main() {
	// Check for required arguments
	if len(os.Args) < 2 {
		fmt.Println("VirtualHere - Service Mode Example")
		fmt.Println("===================================\n")
		fmt.Println("This example demonstrates running VirtualHere client as a managed service")
		fmt.Println("with automatic process monitoring and cleanup.\n")
		fmt.Println("Usage: go run service_mode.go <binary_path>")
		fmt.Println("\nExample:")
		fmt.Println("  go run service_mode.go ./vhclient")
		os.Exit(1)
	}

	binaryPath := os.Args[1]

	fmt.Println("VirtualHere - Service Mode Example")
	fmt.Println("===================================\n")

	// Create client with service mode enabled and process termination callback
	client, err := vh.NewClient(
		binaryPath,
		vh.WithService(true),
		vh.WithOnProcessTerminated(func() {
			fmt.Println("\n⚠️  VirtualHere process terminated externally!")
			fmt.Println("   Resources have been automatically cleaned up.")
			os.Exit(0)
		}),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	fmt.Println("✓ VirtualHere client service started successfully!")
	fmt.Println("\nThe service is now running in the background.")
	fmt.Println("Features:")
	fmt.Println("  • Automatic process monitoring")
	fmt.Println("  • Callback on external termination")
	fmt.Println("  • Automatic resource cleanup")
	fmt.Println("\nTry the following:")
	fmt.Println("  1. List devices: The service is ready to accept commands")
	fmt.Println("  2. Kill the process externally (e.g., pkill vhclient)")
	fmt.Println("  3. Press Ctrl+C for graceful shutdown")

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Periodically check and display device status
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	fmt.Println("\n" + strings.Repeat("-", 50))
	fmt.Println("Monitoring device status (updates every 5 seconds)...")
	fmt.Println(strings.Repeat("-", 50))

	// Initial device list
	displayDevices(client)

	for {
		select {
		case <-ticker.C:
			displayDevices(client)
		case sig := <-sigChan:
			fmt.Printf("\n\nReceived signal: %v\n", sig)
			fmt.Println("Shutting down gracefully...")
			return
		}
	}
}

func displayDevices(client *vh.Client) {
	state, err := client.List()
	if err != nil {
		fmt.Printf("\n[%s] Failed to list devices: %v\n", time.Now().Format("15:04:05"), err)
		return
	}

	fmt.Printf("\n[%s] Connected Hubs: %d\n", time.Now().Format("15:04:05"), len(state.Hubs))

	if len(state.Hubs) == 0 {
		fmt.Println("  No hubs connected. Add a hub with: vh.Client.ManualHub()")
		return
	}

	for _, hub := range state.Hubs {
		fmt.Printf("  • %s (%s) - %d devices\n", hub.Name, hub.Address, len(hub.Devices))
		for _, device := range hub.Devices {
			status := "Available"
			if device.InUse {
				status = "In Use"
			}
			fmt.Printf("    └─ %s (%s) - %s\n", device.Name, device.Address, status)
		}
	}
}
