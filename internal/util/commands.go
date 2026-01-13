package util

import (
	"context"
	"fmt"
	"log"
	"os/exec"

	"github.com/google/shlex"
	"github.com/pkg/errors"
)

type CommandType string

const (
	CommandTypeBash                 CommandType = "bash"
	CommandTypeSh                   CommandType = "sh"
	CommandTypeWindowsCommandPrompt CommandType = "cmd"
	CommandTypeWindowsPowerShell    CommandType = "powershell"
)

var (
	CommandTypeToLaunchTemplate = map[CommandType]string{
		CommandTypeBash:                 "bash -c '%s'",
		CommandTypeSh:                   "sh -c '%s'",
		CommandTypeWindowsCommandPrompt: "cmd /C \"%s\"",
		CommandTypeWindowsPowerShell:    "powershell -Command \"%s\"",
	}
)

func WrapCommand(command string, commandType CommandType) (string, error) {
	t, ok := CommandTypeToLaunchTemplate[commandType]
	if !ok {
		return "", errors.Errorf("unsupported command type: %s", commandType)
	}
	wrappedCommand := fmt.Sprintf(t, command)
	return wrappedCommand, nil
}

func runCommandGetOutput(ctx context.Context, command string) (string, error) {
	argv, err := shlex.Split(command)
	if err != nil {
		return "", errors.Wrap(err, "failed to split command")
	}

	// Resolve the full path to the executable.
	exeName := argv[0]
	exePath, err := exec.LookPath(exeName)
	if err != nil {
		return "", errors.Wrapf(err, "failed to find executable for command: %s", exeName)
	}
	if exeName != exePath {
		log.Printf("Resolved path to executable: %s -> %s", exeName, exePath)
	}

	cmd := exec.Command(exePath, argv[1:]...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func runCommandGetOutputSplitLines(ctx context.Context, command, delimeter string) ([]string, error) {
	output, err := runCommandGetOutput(ctx, command)
	if err != nil {
		return nil, err
	}
	lines := SplitLinesNonEmpty(output, delimeter)
	return lines, nil
}
