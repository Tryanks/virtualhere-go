# virtualhere-go

A Go library for controlling VirtualHere USB clients programmatically.

## Overview

`virtualhere-go` provides a Go interface to interact with VirtualHere USB client software, allowing you to manage USB devices shared over the network.

## Features

- Control VirtualHere USB client from Go applications
- List available USB servers and devices
- Use and stop using remote USB devices
- Manage server connections
- Parse command responses into structured data

## Installation

```bash
go get github.com/Tryanks/virtualhere-go
```

## Usage

```go
package main

import (
    "fmt"
    "log"

    vh "github.com/Tryanks/virtualhere-go"
)

func main() {
    // Create a new client with the path to your VirtualHere binary
    // Windows: vhui64.exe
    // Linux: vhclientx86_64
    // macOS: vhclient
    client, err := vh.NewClient("/path/to/vhclient")
    if err != nil {
        log.Fatal(err)
    }

    // List all available devices
    devices, err := client.ListDevices()
    if err != nil {
        log.Fatal(err)
    }

    for _, device := range devices {
        fmt.Printf("Device: %s (Address: %s)\n", device.Name, device.Address)
    }
}
```

## Requirements

- Go 1.25 or higher
- VirtualHere USB client binary installed on your system
  - Download from: https://www.virtualhere.com/usb_client_software

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Author

Tryanks (tryanks@outlook.com)
