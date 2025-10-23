//go:build windows
// +build windows

package virtualhere

import (
	"fmt"
	"io"
	"time"

	"github.com/Microsoft/go-winio"
)

// executeCommandWindows sends a command via Windows named pipe
func (c *Client) executeCommandWindows(command string) (string, error) {
	// Connect to the named pipe
	pipePath := `\\.\pipe\vhclient`

	timeout := 5 * time.Second
	conn, err := winio.DialPipe(pipePath, &timeout)
	if err != nil {
		return "", fmt.Errorf("failed to connect to named pipe: %w", err)
	}
	defer conn.Close()

	// Set read/write deadline (5 seconds as per VirtualHere API documentation)
	deadline := time.Now().Add(5 * time.Second)
	if err := conn.SetDeadline(deadline); err != nil {
		return "", fmt.Errorf("failed to set deadline: %w", err)
	}

	// Write command with newline (required by VirtualHere protocol)
	// Do NOT include null terminator
	if _, err := conn.Write([]byte(command + "\n")); err != nil {
		return "", fmt.Errorf("failed to write command: %w", err)
	}

	// Read response
	response, err := io.ReadAll(conn)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	return string(response), nil
}

// executeCommandUnix is a stub for Windows (not used)
func (c *Client) executeCommandUnix(command string) (string, error) {
	return "", fmt.Errorf("Unix domain sockets are not supported on Windows")
}
