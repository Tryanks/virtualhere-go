# virtualhere-go

A Go library for controlling VirtualHere USB clients programmatically.

## Overview

`virtualhere-go` provides a Go interface to interact with VirtualHere USB client software, allowing you to manage USB devices shared over the network. It communicates with the VirtualHere client daemon through platform-specific IPC mechanisms:
- **Windows**: Named pipe (`\\.\pipe\vhclient`)
- **Linux/macOS**: Unix domain sockets (`/tmp/vhclient` and `/tmp/vhclient_response`)

## Features

- Control VirtualHere USB client from Go applications
- List available USB servers and devices
- Use and stop using remote USB devices
- Manage server connections and auto-use settings
- Run client as a background service/daemon
- Parse command responses into structured data
- Support for all VirtualHere API commands

## Installation

```bash
go get github.com/Tryanks/virtualhere-go
```

## Requirements

- Go 1.25 or higher
- VirtualHere USB client software installed on your system
  - Download from: https://www.virtualhere.com/usb_client_software
- VirtualHere client must be running as a daemon/service before using this library

## Usage

### Starting the VirtualHere Client Daemon

Before using this library, you need to start the VirtualHere client as a daemon:

**Linux:**
```bash
# Start the daemon
sudo ./vhclientx86_64 -n

# Or use systemd (recommended)
sudo systemctl start virtualhereclient.service
```

**macOS:**
```bash
# Run the client in background
./vhclient &
```

**Windows:**
```powershell
# Install as service (requires admin rights)
vhui64.exe -i

# Or just run the GUI client (it creates the named pipe automatically)
vhui64.exe
```

### Basic Example

```go
package main

import (
    "fmt"
    "log"

    vh "github.com/Tryanks/virtualhere-go"
)

func main() {
    // Create a new client
    // The binary path is only used if you want to manage the service
    // For communication, it uses IPC (named pipe on Windows, Unix sockets on Linux/macOS)
    client, err := vh.NewClient("/path/to/vhclient")
    if err != nil {
        log.Fatal(err)
    }

    // List all available hubs and devices
    state, err := client.List()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Auto-Find: %v\n", state.AutoFindEnabled)
    fmt.Printf("Running as Service: %v\n\n", state.RunningAsService)

    for _, hub := range state.Hubs {
        fmt.Printf("Hub: %s (%s)\n", hub.Name, hub.Address)
        for _, device := range hub.Devices {
            fmt.Printf("  - %s (%s) [In Use: %v]\n",
                device.Name, device.Address, device.InUse)
        }
    }

    // Use a device
    err = client.Use("raspberrypi.114", "")
    if err != nil {
        log.Fatal(err)
    }

    // Stop using a device
    err = client.StopUsing("raspberrypi.114")
    if err != nil {
        log.Fatal(err)
    }
}
```

### Running Client as a Managed Service

You can also let the library manage the VirtualHere client process:

```go
package main

import (
    "fmt"
    "log"

    vh "github.com/Tryanks/virtualhere-go"
)

func main() {
    // Create client with service management
    client, err := vh.NewClient("/path/to/vhclient",
        vh.WithService(true),
        vh.WithOnProcessTerminated(func() {
            fmt.Println("VirtualHere client terminated!")
        }),
    )
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close() // Gracefully shutdown the service

    // Now use the client...
    state, err := client.List()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Found %d hubs\n", len(state.Hubs))
}
```

### More Examples

```go
// Add a manual hub
client.ManualHubAdd("192.168.1.100:7575")

// Enable auto-use for all devices on a hub
client.AutoUseHub("raspberrypi:7575")

// Get device information
info, err := client.DeviceInfo("raspberrypi.114")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Vendor: %s\nProduct: %s\n", info.Vendor, info.Product)

// Get detailed client state as XML
xmlState, err := client.GetClientState()
if err != nil {
    log.Fatal(err)
}
```

## How It Works

This library communicates with the VirtualHere client daemon using platform-specific IPC:

### Windows
Uses Windows Named Pipe at `\\.\pipe\vhclient`. The library opens the pipe, writes the command, and reads the response in message mode.

### Linux/macOS
Uses two Unix domain sockets:
- `/tmp/vhclient` - for sending commands (write only)
- `/tmp/vhclient_response` - for receiving responses (read only)

The command must be terminated with a newline character (`\n`).

## API Documentation

See the [GoDoc](https://pkg.go.dev/github.com/Tryanks/virtualhere-go) for full API documentation.

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Author

Tryanks (tryanks@outlook.com)
