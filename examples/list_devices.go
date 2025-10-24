package main

import (
	"fmt"
	"log"

	vh "github.com/Tryanks/virtualhere-go"
)

func main() {
	// Create VirtualHere pipe client (connects to running service)
	client, err := vh.NewPipeClient()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	fmt.Println("VirtualHere - List Devices Example")
	fmt.Println("===================================\n")

	// List all available devices
	state, err := client.List()
	if err != nil {
		log.Fatalf("Failed to list devices: %v", err)
	}

	// Display client state
	fmt.Printf("Auto-Find: %v\n", state.AutoFindEnabled)
	fmt.Printf("Auto-Use All: %v\n", state.AutoUseAllEnabled)
	fmt.Printf("Reverse Lookup: %v\n", state.ReverseLookup)
	fmt.Printf("Running as service: %v\n\n", state.RunningAsService)

	// Display all hubs and devices
	if len(state.Hubs) == 0 {
		fmt.Println("No hubs found.")
		return
	}

	for _, hub := range state.Hubs {
		fmt.Printf("Hub: %s (%s)\n", hub.Name, hub.Address)

		if len(hub.Devices) == 0 {
			fmt.Println("  (no devices)")
		} else {
			for _, device := range hub.Devices {
				autoUseMarker := ""
				if device.AutoUse {
					autoUseMarker = " *"
				}
				fmt.Printf("  --> %s (%s)%s\n", device.Name, device.Address, autoUseMarker)
			}
		}
		fmt.Println()
	}
}
