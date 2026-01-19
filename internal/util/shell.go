package util

import (
	"runtime"
)

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
