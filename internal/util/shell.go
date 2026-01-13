package util

import (
	"runtime"
)

// GetDefaultShell returns the user's default shell based on the operating system type.
func GetDefaultShell() (string, error) {
	var shell string
	switch runtime.GOOS {
	case "windows":
		shell = GetEnv("COMSPEC", "SHELL")
	case "linux", "darwin":
		shell = GetEnv("SHELL")
	}
	return shell, nil
}
