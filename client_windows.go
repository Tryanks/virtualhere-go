//go:build windows
// +build windows

package virtualhere

import (
	"fmt"
	"io"
	"syscall"
	"time"

	"github.com/Microsoft/go-winio"
)

const (
	PIPE_READMODE_MESSAGE = 0x2
)

var (
	kernel32                = syscall.NewLazyDLL("kernel32.dll")
	setNamedPipeHandleState = kernel32.NewProc("SetNamedPipeHandleState")
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

	// Set pipe to MESSAGE mode (required by VirtualHere)
	// This matches the C++ example from the API docs:
	// DWORD dwMode = PIPE_READMODE_MESSAGE;
	// SetNamedPipeHandleState(hPipe, &dwMode, NULL, NULL)
	if pipeConn, ok := conn.(*winio.PipeConn); ok {
		// Get the underlying file descriptor
		rawConn, err := pipeConn.SyscallConn()
		if err == nil {
			var setErr error
			rawConn.Control(func(fd uintptr) {
				mode := uint32(PIPE_READMODE_MESSAGE)
				ret, _, callErr := setNamedPipeHandleState.Call(
					fd,
					uintptr(syscall.Pointer(&mode)),
					0,
					0,
				)
				if ret == 0 {
					setErr = fmt.Errorf("SetNamedPipeHandleState failed: %w", callErr)
				}
			})
			if setErr != nil {
				return "", setErr
			}
		}
	}

	// Set read/write deadline
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
