package virtualhere

import (
	"fmt"
	"strings"
)

// List returns a list of all available devices and hubs
func (c *Client) List() (*ClientState, error) {
	result, err := c.executeCommand("LIST")
	if err != nil {
		return nil, err
	}

	if !result.Success {
		return nil, result.Error
	}

	return parseListOutput(result.Output)
}

// GetClientState returns the detailed full client state as an XML document
func (c *Client) GetClientState() (*XMLClientState, error) {
	result, err := c.executeCommand("GET CLIENT STATE")
	if err != nil {
		return nil, err
	}

	if !result.Success {
		return nil, result.Error
	}

	return parseClientStateXML(result.Output)
}

// Use connects to and uses a remote device
// address: device address (e.g., "raspberrypi.114")
// password: optional password for the device (empty string if none)
func (c *Client) Use(address string, password string) error {
	var command string
	if password != "" {
		command = fmt.Sprintf("USE,%s,%s", address, password)
	} else {
		command = fmt.Sprintf("USE,%s", address)
	}

	result, err := c.executeCommand(command)
	if err != nil {
		return err
	}

	if !result.Success {
		return result.Error
	}

	return nil
}

// StopUsing disconnects from a device
func (c *Client) StopUsing(address string) error {
	result, err := c.executeCommand(fmt.Sprintf("STOP USING,%s", address))
	if err != nil {
		return err
	}

	if !result.Success {
		return result.Error
	}

	return nil
}

// StopUsingAll stops using all devices on all servers or a specific server
// If serverAddress is empty, stops all devices on all servers
// serverAddress can be in format "address:port" or "EasyFind address"
func (c *Client) StopUsingAll(serverAddress string) error {
	var command string
	if serverAddress != "" {
		command = fmt.Sprintf("STOP USING ALL,%s", serverAddress)
	} else {
		command = "STOP USING ALL"
	}

	result, err := c.executeCommand(command)
	if err != nil {
		return err
	}

	if !result.Success {
		return result.Error
	}

	return nil
}

// StopUsingAllLocal stops using all devices just for this client
func (c *Client) StopUsingAllLocal() error {
	result, err := c.executeCommand("STOP USING ALL LOCAL")
	if err != nil {
		return err
	}

	if !result.Success {
		return result.Error
	}

	return nil
}

// DeviceInfo returns information about a specific device
func (c *Client) DeviceInfo(address string) (*DeviceInfo, error) {
	result, err := c.executeCommand(fmt.Sprintf("DEVICE INFO,%s", address))
	if err != nil {
		return nil, err
	}

	if !result.Success {
		return nil, result.Error
	}

	return parseDeviceInfo(result.Output)
}

// ServerInfo returns information about a specific server
func (c *Client) ServerInfo(serverName string) (*ServerInfo, error) {
	result, err := c.executeCommand(fmt.Sprintf("SERVER INFO,%s", serverName))
	if err != nil {
		return nil, err
	}

	if !result.Success {
		return nil, result.Error
	}

	return parseServerInfo(result.Output)
}

// DeviceRename sets a nickname for a device
func (c *Client) DeviceRename(address string, nickname string) error {
	result, err := c.executeCommand(fmt.Sprintf("DEVICE RENAME,%s,%s", address, nickname))
	if err != nil {
		return err
	}

	if !result.Success {
		return result.Error
	}

	return nil
}

// ServerRename renames a server
func (c *Client) ServerRename(hubAddress string, newName string) error {
	result, err := c.executeCommand(fmt.Sprintf("SERVER RENAME,%s,%s", hubAddress, newName))
	if err != nil {
		return err
	}

	if !result.Success {
		return result.Error
	}

	return nil
}

// AutoUseAll turns auto-use all devices on
func (c *Client) AutoUseAll() error {
	result, err := c.executeCommand("AUTO USE ALL")
	if err != nil {
		return err
	}

	if !result.Success {
		return result.Error
	}

	return nil
}

// AutoUseHub toggles auto-use for all devices on a specific hub
func (c *Client) AutoUseHub(serverName string) error {
	result, err := c.executeCommand(fmt.Sprintf("AUTO USE HUB,%s", serverName))
	if err != nil {
		return err
	}

	if !result.Success {
		return result.Error
	}

	return nil
}

// AutoUsePort toggles auto-use for any device on a specific port
func (c *Client) AutoUsePort(address string) error {
	result, err := c.executeCommand(fmt.Sprintf("AUTO USE PORT,%s", address))
	if err != nil {
		return err
	}

	if !result.Success {
		return result.Error
	}

	return nil
}

// AutoUseDevice toggles auto-use for a specific device on any port
func (c *Client) AutoUseDevice(address string) error {
	result, err := c.executeCommand(fmt.Sprintf("AUTO USE DEVICE,%s", address))
	if err != nil {
		return err
	}

	if !result.Success {
		return result.Error
	}

	return nil
}

// AutoUseDevicePort toggles auto-use for a specific device on a specific port
func (c *Client) AutoUseDevicePort(address string) error {
	result, err := c.executeCommand(fmt.Sprintf("AUTO USE DEVICE PORT,%s", address))
	if err != nil {
		return err
	}

	if !result.Success {
		return result.Error
	}

	return nil
}

// AutoUseClearAll clears all auto-use settings
func (c *Client) AutoUseClearAll() error {
	result, err := c.executeCommand("AUTO USE CLEAR ALL")
	if err != nil {
		return err
	}

	if !result.Success {
		return result.Error
	}

	return nil
}

// ManualHubAdd adds a manually specified hub to connect to
// address can be in format "address:port" or "EasyFind address"
func (c *Client) ManualHubAdd(address string) error {
	result, err := c.executeCommand(fmt.Sprintf("MANUAL HUB ADD,%s", address))
	if err != nil {
		return err
	}

	if !result.Success {
		return result.Error
	}

	return nil
}

// ManualHubRemove removes a manually specified hub
func (c *Client) ManualHubRemove(address string) error {
	result, err := c.executeCommand(fmt.Sprintf("MANUAL HUB REMOVE,%s", address))
	if err != nil {
		return err
	}

	if !result.Success {
		return result.Error
	}

	return nil
}

// ManualHubRemoveAll removes all manually specified hubs
func (c *Client) ManualHubRemoveAll() error {
	result, err := c.executeCommand("MANUAL HUB REMOVE ALL")
	if err != nil {
		return err
	}

	if !result.Success {
		return result.Error
	}

	return nil
}

// ManualHubList returns a list of manually specified hubs
func (c *Client) ManualHubList() ([]string, error) {
	result, err := c.executeCommand("MANUAL HUB LIST")
	if err != nil {
		return nil, err
	}

	if !result.Success {
		return nil, result.Error
	}

	// Parse the list of hubs from the output
	hubs := make([]string, 0)
	lines := strings.Split(result.Output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "Manual") {
			hubs = append(hubs, line)
		}
	}

	return hubs, nil
}

// AddReverse adds a reverse client to the server
func (c *Client) AddReverse(serverSerial string, clientAddress string) error {
	result, err := c.executeCommand(fmt.Sprintf("ADD REVERSE,%s,%s", serverSerial, clientAddress))
	if err != nil {
		return err
	}

	if !result.Success {
		return result.Error
	}

	return nil
}

// RemoveReverse removes a reverse client from the server
func (c *Client) RemoveReverse(serverSerial string, clientAddress string) error {
	result, err := c.executeCommand(fmt.Sprintf("REMOVE REVERSE,%s,%s", serverSerial, clientAddress))
	if err != nil {
		return err
	}

	if !result.Success {
		return result.Error
	}

	return nil
}

// ListReverse lists all reverse clients for a server
func (c *Client) ListReverse(serverSerial string) ([]string, error) {
	result, err := c.executeCommand(fmt.Sprintf("LIST REVERSE,%s", serverSerial))
	if err != nil {
		return nil, err
	}

	if !result.Success {
		return nil, result.Error
	}

	// Parse the list of reverse clients
	clients := make([]string, 0)
	lines := strings.Split(result.Output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			clients = append(clients, line)
		}
	}

	return clients, nil
}

// ListLicenses returns a list of licenses
func (c *Client) ListLicenses() ([]string, error) {
	result, err := c.executeCommand("LIST LICENSES")
	if err != nil {
		return nil, err
	}

	if !result.Success {
		return nil, result.Error
	}

	// Parse the list of licenses
	licenses := make([]string, 0)
	lines := strings.Split(result.Output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			licenses = append(licenses, line)
		}
	}

	return licenses, nil
}

// LicenseServer licenses a server with a license key
func (c *Client) LicenseServer(licenseKey string) error {
	result, err := c.executeCommand(fmt.Sprintf("LICENSE SERVER,%s", licenseKey))
	if err != nil {
		return err
	}

	if !result.Success {
		return result.Error
	}

	return nil
}

// ClearLog clears the client log
func (c *Client) ClearLog() error {
	result, err := c.executeCommand("CLEAR LOG")
	if err != nil {
		return err
	}

	if !result.Success {
		return result.Error
	}

	return nil
}

// CustomEvent sets a custom device event
func (c *Client) CustomEvent(address string, event string) error {
	result, err := c.executeCommand(fmt.Sprintf("CUSTOM EVENT,%s,%s", address, event))
	if err != nil {
		return err
	}

	if !result.Success {
		return result.Error
	}

	return nil
}

// AutoFind toggles auto-find functionality
func (c *Client) AutoFind() error {
	result, err := c.executeCommand("AUTOFIND")
	if err != nil {
		return err
	}

	if !result.Success {
		return result.Error
	}

	return nil
}

// Reverse toggles reverse lookup functionality
func (c *Client) Reverse() error {
	result, err := c.executeCommand("REVERSE")
	if err != nil {
		return err
	}

	if !result.Success {
		return result.Error
	}

	return nil
}

// SSLReverse toggles reverse SSL lookup
func (c *Client) SSLReverse() error {
	result, err := c.executeCommand("SSLREVERSE")
	if err != nil {
		return err
	}

	if !result.Success {
		return result.Error
	}

	return nil
}

// Exit shuts down the client
func (c *Client) Exit() error {
	result, err := c.executeCommand("EXIT")
	if err != nil {
		return err
	}

	if !result.Success {
		return result.Error
	}

	return nil
}

// Help returns the help message with available commands
func (c *Client) Help() (string, error) {
	result, err := c.executeCommand("HELP")
	if err != nil {
		return "", err
	}

	// Help command always returns output even on failure
	return result.Output, nil
}
