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
		fmt.Println("VirtualHere - Auto-Use Configuration Example")
		fmt.Println("=============================================\n")
		fmt.Println("Usage: go run auto_use.go <command> [arguments]")
		fmt.Println("\nCommands:")
		fmt.Println("  hub <hub_address>      - Toggle auto-use for all devices on a hub")
		fmt.Println("  device <device_addr>   - Toggle auto-use for a device on any port")
		fmt.Println("  port <device_addr>     - Toggle auto-use for any device on a port")
		fmt.Println("  all                    - Enable auto-use for all devices")
		fmt.Println("  clear                  - Clear all auto-use settings")
		fmt.Println("\nExamples:")
		fmt.Println("  go run auto_use.go hub raspberrypi:7575")
		fmt.Println("  go run auto_use.go device raspberrypi.114")
		fmt.Println("  go run auto_use.go all")
		os.Exit(1)
	}

	command := os.Args[1]

	// Create VirtualHere pipe client (connects to running service)
	client, err := vh.NewPipeClient()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	fmt.Println("VirtualHere - Auto-Use Configuration")
	fmt.Println("====================================\n")

	switch command {
	case "hub":
		if len(os.Args) < 3 {
			log.Fatal("Hub address is required. Example: raspberrypi:7575")
		}
		hubAddress := os.Args[2]
		fmt.Printf("Toggling auto-use for hub: %s\n", hubAddress)
		err = client.AutoUseHub(hubAddress)
		if err != nil {
			log.Fatalf("Failed to toggle auto-use: %v", err)
		}
		fmt.Println("✓ Auto-use toggled successfully!")

	case "device":
		if len(os.Args) < 3 {
			log.Fatal("Device address is required. Example: raspberrypi.114")
		}
		deviceAddress := os.Args[2]
		fmt.Printf("Toggling auto-use for device: %s (any port)\n", deviceAddress)
		err = client.AutoUseDevice(deviceAddress)
		if err != nil {
			log.Fatalf("Failed to toggle auto-use: %v", err)
		}
		fmt.Println("✓ Auto-use toggled successfully!")

	case "port":
		if len(os.Args) < 3 {
			log.Fatal("Device address is required. Example: raspberrypi.114")
		}
		deviceAddress := os.Args[2]
		fmt.Printf("Toggling auto-use for port: %s (any device)\n", deviceAddress)
		err = client.AutoUsePort(deviceAddress)
		if err != nil {
			log.Fatalf("Failed to toggle auto-use: %v", err)
		}
		fmt.Println("✓ Auto-use toggled successfully!")

	case "all":
		fmt.Println("Enabling auto-use for all devices...")
		err = client.AutoUseAll()
		if err != nil {
			log.Fatalf("Failed to enable auto-use all: %v", err)
		}
		fmt.Println("✓ Auto-use all enabled successfully!")

	case "clear":
		fmt.Println("Clearing all auto-use settings...")
		err = client.AutoUseClearAll()
		if err != nil {
			log.Fatalf("Failed to clear auto-use settings: %v", err)
		}
		fmt.Println("✓ All auto-use settings cleared successfully!")

	default:
		log.Fatalf("Unknown command: %s", command)
	}

	// Show current state
	fmt.Println("\nCurrent device state:")
	state, err := client.List()
	if err != nil {
		log.Printf("Failed to list devices: %v", err)
		return
	}

	fmt.Printf("Auto-Use All: %v\n\n", state.AutoUseAllEnabled)
	for _, hub := range state.Hubs {
		fmt.Printf("Hub: %s (%s)\n", hub.Name, hub.Address)
		for _, device := range hub.Devices {
			autoUseMarker := ""
			if device.AutoUse {
				autoUseMarker = " [AUTO-USE ENABLED]"
			}
			fmt.Printf("  --> %s (%s)%s\n", device.Name, device.Address, autoUseMarker)
		}
	}
}
