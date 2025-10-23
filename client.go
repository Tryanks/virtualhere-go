package virtualhere

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Client represents a VirtualHere USB client controller
type Client struct {
	binaryPath string
}

// NewClient creates a new VirtualHere client controller with the specified binary path
// The binary should be the path to vhui64.exe (Windows), vhclientx86_64 (Linux), or vhclient (macOS)
func NewClient(binaryPath string) (*Client, error) {
	// Verify the binary exists
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("%w: %s", ErrBinaryNotFound, binaryPath)
	}

	// Verify the binary is executable (Unix-like systems)
	info, err := os.Stat(binaryPath)
	if err != nil {
		return nil, err
	}

	// Check if file has execute permissions (on Unix-like systems)
	if info.Mode()&0111 == 0 {
		// Try to make it executable
		if err := os.Chmod(binaryPath, 0755); err != nil {
			return nil, fmt.Errorf("binary is not executable and cannot be made executable: %w", err)
		}
	}

	return &Client{
		binaryPath: binaryPath,
	}, nil
}

// executeCommand executes a VirtualHere command using the -t flag
// It returns the raw output and any error that occurred
func (c *Client) executeCommand(command string) (*CommandResult, error) {
	cmd := exec.Command(c.binaryPath, "-t", command)

	output, err := cmd.CombinedOutput()
	outputStr := strings.TrimSpace(string(output))

	result := &CommandResult{
		Output: outputStr,
	}

	// Check the exit code
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			switch exitErr.ExitCode() {
			case 1:
				// FAILED - command failed (e.g., trying to use an in-use device)
				result.Success = false
				result.Error = ErrCommandFailed
				return result, nil
			case 2:
				// ERROR - error occurred (e.g., server doesn't exist, invalid address)
				result.Success = false
				// Try to extract the error message
				if strings.HasPrefix(outputStr, "ERROR:") {
					result.Error = fmt.Errorf("%s", strings.TrimPrefix(outputStr, "ERROR:"))
				} else {
					result.Error = fmt.Errorf("%s", outputStr)
				}
				return result, nil
			default:
				return result, fmt.Errorf("command execution failed with exit code %d: %w", exitErr.ExitCode(), err)
			}
		}
		return result, fmt.Errorf("command execution failed: %w", err)
	}

	// Exit code 0 means OK
	result.Success = true
	return result, nil
}

// GetBinaryPath returns the path to the VirtualHere binary
func (c *Client) GetBinaryPath() string {
	return c.binaryPath
}
