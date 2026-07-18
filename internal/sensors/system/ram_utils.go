package system

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"os/exec"
)

// executeCommand executes a command and returns its output
func executeCommand(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	return cmd.Output()
}

// Update method for FreeRAMSensor that dispatches to the appropriate OS-specific implementation
func (frs *FreeRAMSensor) Update(ctx context.Context) (err error) {
	switch runtime.GOOS {
	case "darwin":
		return frs.updateDarwin(ctx)
	case "windows":
		return frs.updateWindows(ctx)
	case "linux":
		// Fall through to the original Linux implementation
		return frs.updateLinux(ctx)
	default:
		frs.Mutex.Lock()
		defer frs.Mutex.Unlock()
		frs.UpdatedAt = time.Now()
		err = fmt.Errorf("RAM sensor not implemented on %s", runtime.GOOS)
		frs.Error = err
		return err
	}
}

// updateLinux contains the original Linux implementation
// It's renamed to avoid conflict with the new Update method
func (frs *FreeRAMSensor) updateLinux(ctx context.Context) (err error) {
	// This method will be implemented with the original Linux code
	// It's declared here to avoid conflicts with the new Update method
	// The actual implementation remains in ram.go but will need to be adjusted
	return
}
