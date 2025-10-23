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
	binaryPath           string
	serviceCmd           *exec.Cmd
	serviceMu            sync.Mutex
	runService           bool
	onProcessTerminated  func()
	processMonitorDone   chan struct{}
	processMonitorCancel chan struct{}
}

// ClientOption is a function that configures a Client
type ClientOption func(*Client)

// WithService configures the client to run as a background service
func WithService(enable bool) ClientOption {
	return func(c *Client) {
		c.runService = enable
	}
}

// WithOnProcessTerminated sets a callback function that will be called when the
// managed service process is terminated externally or by the user.
// The client will automatically cleanup resources when this occurs.
// This option only works when WithService(true) is enabled.
func WithOnProcessTerminated(callback func()) ClientOption {
	return func(c *Client) {
		c.onProcessTerminated = callback
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

	// Start monitoring the process for termination
	c.processMonitorCancel = make(chan struct{})
	c.processMonitorDone = make(chan struct{})
	go c.monitorProcess()

	return nil
}

// monitorProcess monitors the service process for termination and triggers cleanup
func (c *Client) monitorProcess() {
	defer close(c.processMonitorDone)

	// Wait for process to exit or cancellation
	processDone := make(chan error, 1)
	go func() {
		processDone <- c.serviceCmd.Wait()
	}()

	select {
	case <-processDone:
		// Process terminated externally or by EXIT command
		c.serviceMu.Lock()
		c.serviceCmd = nil
		callback := c.onProcessTerminated
		c.serviceMu.Unlock()

		// Call the termination callback if set
		if callback != nil {
			callback()
		}
	case <-c.processMonitorCancel:
		// Close() was called, normal shutdown
		return
	}
}

// Close stops the background service if running
func (c *Client) Close() error {
	c.serviceMu.Lock()
	defer c.serviceMu.Unlock()

	if c.serviceCmd == nil {
		return nil
	}

	// Signal the monitor to stop
	if c.processMonitorCancel != nil {
		close(c.processMonitorCancel)
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

	// Wait for monitor goroutine to finish
	if c.processMonitorDone != nil {
		<-c.processMonitorDone
	}

	return nil
}

// executeCommand sends a command to the VirtualHere client via named pipe (Windows)
// or Unix socket (Linux/macOS) and returns the response
func (c *Client) executeCommand(command string) (*CommandResult, error) {
	result := &CommandResult{}

	var response string
	var err error

	if runtime.GOOS == "windows" {
		response, err = c.executeCommandWindows(command)
	} else {
		response, err = c.executeCommandUnix(command)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to communicate with client: %w", err)
	}

	// Parse the response
	response = strings.TrimSpace(response)
	result.Output = response

	// Check response status
	if response == "OK" {
		result.Success = true
		return result, nil
	}

	if response == "FAILED" {
		result.Success = false
		result.Error = ErrCommandFailed
		return result, nil
	}

	if strings.HasPrefix(response, "ERROR:") {
		result.Success = false
		result.Error = fmt.Errorf("%s", strings.TrimSpace(strings.TrimPrefix(response, "ERROR:")))
		return result, nil
	}

	// If response is not OK/FAILED/ERROR, it's likely data (e.g., from LIST command)
	result.Success = true
	return result, nil
}

// GetBinaryPath returns the path to the VirtualHere binary
func (c *Client) GetBinaryPath() string {
	return c.binaryPath
}
