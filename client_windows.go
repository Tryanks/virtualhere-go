//go:build windows
// +build windows

package virtualhere

import (
	"fmt"
	"syscall"
	"time"
	"unsafe"
)

const (
	PIPE_READMODE_MESSAGE = 0x2
	GENERIC_READ          = 0x80000000
	GENERIC_WRITE         = 0x40000000
	OPEN_EXISTING         = 3
	INVALID_HANDLE_VALUE  = ^uintptr(0)
	ERROR_PIPE_BUSY       = 231
)

var (
	kernel32                = syscall.NewLazyDLL("kernel32.dll")
	createFileW             = kernel32.NewProc("CreateFileW")
	setNamedPipeHandleState = kernel32.NewProc("SetNamedPipeHandleState")
	waitNamedPipeW          = kernel32.NewProc("WaitNamedPipeW")
	readFile                = kernel32.NewProc("ReadFile")
	writeFile               = kernel32.NewProc("WriteFile")
	closeHandle             = kernel32.NewProc("CloseHandle")
)

// executeCommandWindows sends a command via Windows named pipe
func (c *Client) executeCommandWindows(command string) (string, error) {
	pipePath := `\\.\pipe\vhclient`
	pipePathPtr, err := syscall.UTF16PtrFromString(pipePath)
	if err != nil {
		return "", fmt.Errorf("failed to convert pipe path: %w", err)
	}

	// Try to open the pipe
	var hPipe uintptr
	timeout := time.Now().Add(5 * time.Second)

	for {
		handle, _, _ := createFileW.Call(
			uintptr(unsafe.Pointer(pipePathPtr)),
			GENERIC_READ|GENERIC_WRITE,
			0,
			0,
			OPEN_EXISTING,
			0,
			0,
		)

		if handle != INVALID_HANDLE_VALUE {
			hPipe = handle
			break
		}

		// Check if pipe is busy
		lastErr := syscall.GetLastError()
		if lastErr == syscall.Errno(ERROR_PIPE_BUSY) {
			// Wait for pipe to become available (2 seconds timeout)
			ret, _, _ := waitNamedPipeW.Call(
				uintptr(unsafe.Pointer(pipePathPtr)),
				2000,
			)
			if ret == 0 {
				return "", fmt.Errorf("failed to wait for named pipe: timeout")
			}
			// Retry opening
			if time.Now().After(timeout) {
				return "", fmt.Errorf("failed to connect to named pipe: timeout")
			}
			continue
		}

		return "", fmt.Errorf("failed to open named pipe: %v", lastErr)
	}

	defer closeHandle.Call(hPipe)

	// Set pipe to MESSAGE mode (required by VirtualHere)
	// This matches the C++ example from the API docs:
	// DWORD dwMode = PIPE_READMODE_MESSAGE;
	// SetNamedPipeHandleState(hPipe, &dwMode, NULL, NULL)
	mode := uint32(PIPE_READMODE_MESSAGE)
	ret, _, callErr := setNamedPipeHandleState.Call(
		hPipe,
		uintptr(unsafe.Pointer(&mode)),
		0,
		0,
	)
	if ret == 0 {
		return "", fmt.Errorf("failed to set pipe mode: %v", callErr)
	}

	// Write command with newline (required by VirtualHere protocol)
	// Do NOT include null terminator
	cmdBytes := []byte(command + "\n")
	var bytesWritten uint32
	ret, _, callErr = writeFile.Call(
		hPipe,
		uintptr(unsafe.Pointer(&cmdBytes[0])),
		uintptr(len(cmdBytes)),
		uintptr(unsafe.Pointer(&bytesWritten)),
		0,
	)
	if ret == 0 {
		return "", fmt.Errorf("failed to write to pipe: %v", callErr)
	}

	// Read response
	buffer := make([]byte, 8192)
	var bytesRead uint32
	ret, _, callErr = readFile.Call(
		hPipe,
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(len(buffer)),
		uintptr(unsafe.Pointer(&bytesRead)),
		0,
	)
	if ret == 0 {
		return "", fmt.Errorf("failed to read from pipe: %v", callErr)
	}

	return string(buffer[:bytesRead]), nil
}

// executeCommandUnix is a stub for Windows (not used)
func (c *Client) executeCommandUnix(command string) (string, error) {
	return "", fmt.Errorf("Unix domain sockets are not supported on Windows")
}
