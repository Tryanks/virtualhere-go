package virtualhere

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

// Client represents a VirtualHere USB client controller
type Client struct {
	binaryPath string
	serviceCmd *exec.Cmd
	serviceMu  sync.Mutex
	runService bool
}

// ClientOption is a function that configures a Client
type ClientOption func(*Client)

// WithService configures the client to run as a background service
func WithService(enable bool) ClientOption {
	return func(c *Client) {
		c.runService = enable
	}
}

// NewClient creates a new VirtualHere client controller with the specified binary path
// The binary should be the path to vhui64.exe (Windows), vhclientx86_64 (Linux), or vhclient (macOS)
// Options can be provided to configure the client behavior
func NewClient(binaryPath string, opts ...ClientOption) (*Client, error) {
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

	client := &Client{
		binaryPath: binaryPath,
	}

	// Apply options
	for _, opt := range opts {
		opt(client)
	}

	// Start service if enabled
	if client.runService {
		if err := client.startService(); err != nil {
			return nil, fmt.Errorf("failed to start service: %w", err)
		}
	}

	return client, nil
}

// startService starts the VirtualHere client as a background service
func (c *Client) startService() error {
	c.serviceMu.Lock()
	defer c.serviceMu.Unlock()

	if c.serviceCmd != nil {
		return fmt.Errorf("service is already running")
	}

	// Start the client based on platform:
	// - Linux: Use -n flag for daemon mode (Console Client)
	// - Windows/macOS: Start without parameters (GUI Client runs in background)
	if runtime.GOOS == "linux" {
		// Linux Console Client uses -n for daemon mode
		c.serviceCmd = exec.Command(c.binaryPath, "-n")
	} else {
		// Windows/macOS GUI clients run in background without -n
		c.serviceCmd = exec.Command(c.binaryPath)
	}

	if err := c.serviceCmd.Start(); err != nil {
		c.serviceCmd = nil
		return fmt.Errorf("failed to start service: %w", err)
	}

	return nil
}

// Close stops the background service if running
func (c *Client) Close() error {
	c.serviceMu.Lock()
	defer c.serviceMu.Unlock()

	if c.serviceCmd == nil {
		return nil
	}

	// Try graceful shutdown first using EXIT command
	if err := c.Exit(); err != nil {
		// If EXIT command fails, force kill the process
		if killErr := c.serviceCmd.Process.Kill(); killErr != nil {
			return fmt.Errorf("failed to kill service process: %w", killErr)
		}
	}

	// Wait for the process to exit
	_ = c.serviceCmd.Wait()
	c.serviceCmd = nil

	return nil
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
