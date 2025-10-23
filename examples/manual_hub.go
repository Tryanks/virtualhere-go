package main

import (
	"fmt"
	"log"
	"os"

	vh "github.com/Tryanks/virtualhere-go"
)

func main() {
	// Check for required arguments
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run manual_hub.go <binary_path> <hub_address>")
		fmt.Println("Example: go run manual_hub.go ./vhclient 192.168.1.100:7575")
		fmt.Println("Example: go run manual_hub.go ./vhclient myserver.easyfind.com")
		os.Exit(1)
	}

	binaryPath := os.Args[1]
	hubAddress := os.Args[2]

	// Create VirtualHere client
	client, err := vh.NewClient(binaryPath)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	fmt.Println("VirtualHere - Manual Hub Management Example")
	fmt.Println("============================================\n")

	// Add manual hub
	fmt.Printf("Adding hub: %s\n", hubAddress)
	err = client.ManualHubAdd(hubAddress)
	if err != nil {
		log.Fatalf("Failed to add hub: %v", err)
	}
	fmt.Println("✓ Hub added successfully!")

	// List manual hubs
	fmt.Println("\nListing all manual hubs...")
	hubs, err := client.ManualHubList()
	if err != nil {
		log.Fatalf("Failed to list hubs: %v", err)
	}

	if len(hubs) == 0 {
		fmt.Println("No manual hubs configured.")
	} else {
		for i, hub := range hubs {
			fmt.Printf("%d. %s\n", i+1, hub)
		}
	}

	// Wait for user input before removing
	fmt.Println("\nPress Enter to remove the hub...")
	fmt.Scanln()

	// Remove the hub
	fmt.Printf("Removing hub: %s\n", hubAddress)
	err = client.ManualHubRemove(hubAddress)
	if err != nil {
		log.Fatalf("Failed to remove hub: %v", err)
	}
	fmt.Println("✓ Hub removed successfully!")
}
