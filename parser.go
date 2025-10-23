package virtualhere

import (
	"encoding/xml"
	"strings"
)

// parseListOutput parses the output of the LIST command
// It extracts hubs and devices from the formatted text output
func parseListOutput(output string) (*ClientState, error) {
	state := &ClientState{
		Hubs: make([]Hub, 0),
	}

	lines := strings.Split(output, "\n")
	var currentHub *Hub

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip header and empty lines
		if line == "" || strings.HasPrefix(line, "VirtualHere IPC") ||
			strings.HasPrefix(line, "(Value in brackets") {
			continue
		}

		// Parse hub line (no leading arrow)
		if !strings.HasPrefix(line, "-->") && strings.Contains(line, "(") && strings.Contains(line, ")") {
			hub := parseHubLine(line)
			if hub.Name != "" {
				state.Hubs = append(state.Hubs, hub)
				currentHub = &state.Hubs[len(state.Hubs)-1]
			}
			continue
		}

		// Parse device line (starts with -->)
		if strings.HasPrefix(line, "-->") && currentHub != nil {
			device := parseDeviceLine(line)
			if device.Name != "" {
				currentHub.Devices = append(currentHub.Devices, device)
			}
			continue
		}

		// Parse status lines
		if strings.Contains(line, "Auto-Find currently") {
			state.AutoFindEnabled = strings.Contains(line, "on")
		}
		if strings.Contains(line, "Auto-Use All currently") {
			state.AutoUseAllEnabled = strings.Contains(line, "on")
		}
		if strings.Contains(line, "Reverse Lookup currently") {
			state.ReverseLookup = strings.Contains(line, "on")
		}
		if strings.Contains(line, "running as a service") {
			state.RunningAsService = !strings.Contains(line, "not")
		}
	}

	return state, nil
}

// parseHubLine parses a hub line from the LIST output
// Format: "Hub Name (address:port)"
func parseHubLine(line string) Hub {
	hub := Hub{
		Devices: make([]Device, 0),
	}

	// Find the address in parentheses
	startIdx := strings.LastIndex(line, "(")
	endIdx := strings.LastIndex(line, ")")

	if startIdx == -1 || endIdx == -1 || startIdx >= endIdx {
		return hub
	}

	hub.Address = strings.TrimSpace(line[startIdx+1 : endIdx])
	hub.Name = strings.TrimSpace(line[:startIdx])

	return hub
}

// parseDeviceLine parses a device line from the LIST output
// Format: "--> Device Name (address)" or "--> Device Name (address) *" (if auto-use)
func parseDeviceLine(line string) Device {
	device := Device{}

	// Remove the arrow prefix
	line = strings.TrimPrefix(line, "-->")
	line = strings.TrimSpace(line)

	// Check for auto-use marker (*)
	if strings.HasSuffix(line, "*") {
		device.AutoUse = true
		line = strings.TrimSuffix(line, "*")
		line = strings.TrimSpace(line)
	}

	// Find the address in parentheses
	startIdx := strings.LastIndex(line, "(")
	endIdx := strings.LastIndex(line, ")")

	if startIdx == -1 || endIdx == -1 || startIdx >= endIdx {
		return device
	}

	device.Address = strings.TrimSpace(line[startIdx+1 : endIdx])
	device.Name = strings.TrimSpace(line[:startIdx])

	return device
}

// parseClientStateXML parses the XML output from GET CLIENT STATE command
func parseClientStateXML(output string) (*XMLClientState, error) {
	var state XMLClientState
	err := xml.Unmarshal([]byte(output), &state)
	if err != nil {
		return nil, err
	}
	return &state, nil
}

// parseDeviceInfo parses the output of DEVICE INFO command
// Example output:
// ADDRESS: TryanksPC.14
// VENDOR: Xiaomi
// VENDOR ID: 0x18d1
// PRODUCT: Mi 10
// PRODUCT ID: 0x4ee7
// SERIAL: 652e1e0d
// IN USE BY: NO ONE
func parseDeviceInfo(output string) (*DeviceInfo, error) {
	info := &DeviceInfo{}
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "ADDRESS":
			info.Address = value
		case "VENDOR":
			info.Vendor = value
		case "VENDOR ID":
			info.VendorID = value
		case "PRODUCT":
			info.Product = value
		case "PRODUCT ID":
			info.ProductID = value
		case "SERIAL":
			info.Serial = value
		case "IN USE BY":
			info.InUseBy = value
		}
	}

	return info, nil
}

// parseServerInfo parses the output of SERVER INFO command
// Example output:
// NAME: Windows Hub
// VERSION: 4.6.4
// STATE: Logged in
// ADDRESS: 192.168.31.145 (192.168.31.145)
// PORT: 7575
// CONNECTED FOR: 9265 sec
// MAX DEVICES: 1
// CONNECTION ID: 1
// INTERFACE:
// SERIAL NUMBER: 07370b72-f03f-4f6e-b930-33fd5d8930f5
// EASYFIND: not enabled
func parseServerInfo(output string) (*ServerInfo, error) {
	info := &ServerInfo{}
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "NAME":
			info.Name = value
		case "VERSION":
			info.Version = value
		case "STATE":
			info.State = value
		case "ADDRESS":
			info.Address = value
		case "PORT":
			info.Port = value
		case "CONNECTED FOR":
			info.ConnectedFor = value
		case "MAX DEVICES":
			info.MaxDevices = value
		case "CONNECTION ID":
			info.ConnectionID = value
		case "INTERFACE":
			info.Interface = value
		case "SERIAL NUMBER":
			info.SerialNumber = value
		case "EASYFIND":
			info.EasyFind = value
		}
	}

	return info, nil
}

// isSuccessResponse checks if the response indicates success
func isSuccessResponse(output string) bool {
	return strings.HasPrefix(output, "OK")
}

// isFailedResponse checks if the response indicates failure
func isFailedResponse(output string) bool {
	return strings.HasPrefix(output, "FAILED")
}

// isErrorResponse checks if the response indicates an error
func isErrorResponse(output string) bool {
	return strings.HasPrefix(output, "ERROR")
}
