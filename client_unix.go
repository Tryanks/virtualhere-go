//go:build !windows
// +build !windows

package virtualhere

import (
	"fmt"
	"io"
	"net"
	"time"
)

// executeCommandWindows is a stub for Unix systems (not used)
func (c *Client) executeCommandWindows(command string) (string, error) {
	return "", fmt.Errorf("Windows named pipe is not supported on this platform")
}

// executeCommandUnix sends a command via Unix domain socket (Linux/macOS)
// The client uses two separate socket files:
// - /tmp/vhclient for sending requests
// - /tmp/vhclient_response for receiving responses
func (c *Client) executeCommandUnix(command string) (string, error) {
	requestPath := "/tmp/vhclient"
	responsePath := "/tmp/vhclient_response"

	// Open response socket first and wait for data
	responseConn, err := net.DialTimeout("unix", responsePath, 2*time.Second)
	if err != nil {
		return "", fmt.Errorf("failed to connect to response socket: %w", err)
	}
	defer responseConn.Close()

	// Connect to request socket
	conn, err := net.DialTimeout("unix", requestPath, 2*time.Second)
	if err != nil {
		return "", fmt.Errorf("failed to connect to request socket: %w", err)
	}
	defer conn.Close()

	// Set write deadline
	if err := conn.SetWriteDeadline(time.Now().Add(2 * time.Second)); err != nil {
		return "", fmt.Errorf("failed to set write deadline: %w", err)
	}

	// Write command with newline (required by VirtualHere protocol)
	if _, err := conn.Write([]byte(command + "\n")); err != nil {
		return "", fmt.Errorf("failed to write command: %w", err)
	}

	// Set read deadline (5 seconds as per VirtualHere API documentation)
	if err := responseConn.SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
		return "", fmt.Errorf("failed to set read deadline: %w", err)
	}

	// Read response - VirtualHere sends complete response
	response, err := io.ReadAll(responseConn)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	return string(response), nil
}
