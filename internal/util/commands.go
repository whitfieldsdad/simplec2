package util

import (
	"context"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/google/shlex"
	"github.com/hairyhenderson/go-which"
	"github.com/pkg/errors"
)

type CommandType string

const (
	CommandTypeNative               CommandType = "native"
	CommandTypeDefaultShell         CommandType = "default-shell"
	CommandTypeWindowsPowerShell    CommandType = "windows-powershell"
	CommandTypeWindowsCommandPrompt CommandType = "windows-command-prompt"
	CommandTypePwsh                 CommandType = "pwsh"
	CommandTypeBash                 CommandType = "bash"
	CommandTypeSh                   CommandType = "sh"
)

var (
	launchTemplates = map[CommandType][]string{
		CommandTypeWindowsPowerShell:    {"powershell", "-ExecutionPolicy", "Bypass", "-Command", "%s"},
		CommandTypeWindowsCommandPrompt: {"cmd", "/c", "%s"},
		CommandTypePwsh:                 {"pwsh", "-Command", "%s"},
		CommandTypeBash:                 {"bash", "-c", "%s"},
		CommandTypeSh:                   {"sh", "-c", "%s"},
	}
)

func ExecuteCommand(ctx context.Context, command string, outputFilePath string) (*Process, error) {
	argv, err := shlex.Split(command)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse command")
	}
	_, err = ExecuteArgv(ctx, argv, outputFilePath)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func ExecuteArgv(ctx context.Context, argv []string, outputFilePath string) (*Process, error) {
	outputFile, err := os.OpenFile(outputFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open output file")
	}
	defer outputFile.Close()

	path, err := getExecutablePath(argv[0])
	if err != nil {
		return nil, errors.Wrap(err, "failed to find executable")
	}

	cmd := exec.CommandContext(ctx, path, argv[1:]...)
	cmd.SysProcAttr = getSysProcAttrs()

	cmd.Stdout = outputFile
	cmd.Stderr = outputFile

	startTime := time.Now()
	err = cmd.Start()
	if err != nil {
		return nil, errors.Wrap(err, "failed to start command")
	}
	pid := cmd.Process.Pid
	ppid := os.Getpid()

	log.Infof("Executing command: `%s` (PID: %d, PPID: %d)", strings.Join(argv, " "), pid, ppid)

	// Gather information about the subprocess.
	process := &Process{
		PID:         cmd.Process.Pid,
		PPID:        os.Getpid(),
		Executable:  path,
		CommandLine: strings.Join(argv, " "),
		StartTime:   &startTime,
	}

	// Wait for the command to complete.
	err = cmd.Wait()
	if err != nil {
		return nil, errors.Wrap(err, "failed to wait for command")
	}

	// Set the exit code and exit time.
	exitCode := cmd.ProcessState.ExitCode()
	process.ExitCode = &exitCode

	endTime := time.Now()
	process.ExitTime = &endTime

	if exitCode == 0 {
		log.Infof("Command exited: `%s` (PID: %d, PPID: %d, exit code: %d)", strings.Join(argv, " "), pid, ppid, exitCode)
	} else {
		log.Warnf("Command exited with non-zero exit code: `%s` (PID: %d, PPID: %d, exit code: %d)", strings.Join(argv, " "), pid, ppid, exitCode)
	}

	return process, nil
}

func getExecutablePath(command string) (string, error) {
	path := which.Which(command)
	if path == "" {
		return "", errors.Errorf("executable for command not found: %s", command)
	}
	return path, nil
}

func WrapCommand(command string, commandType CommandType) ([]string, error) {
	template, ok := launchTemplates[commandType]
	if !ok {
		return nil, errors.New("unsupported command type")
	}
	var wrappedCommand []string
	for _, part := range template {
		if part == "%s" {
			wrappedCommand = append(wrappedCommand, command)
			break
		}
	}
	return wrappedCommand, nil
}
