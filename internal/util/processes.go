package util

import (
	"time"

	"github.com/shirou/gopsutil/v4/process"
)

type Process struct {
	PID         int        `json:"pid"`
	PPID        int        `json:"ppid"`
	Executable  string     `json:"executable,omitempty"`
	CommandLine string     `json:"command_line,omitempty"`
	StartTime   *time.Time `json:"start_time,omitempty"`
	ExitTime    *time.Time `json:"exit_time,omitempty"`
	ExitCode    *int       `json:"exit_code,omitempty"`
}

func (o Process) GetArtifactType() ArtifactType {
	return ArtifactTypeProcess
}

func GetProcess(pid int) (*Process, error) {
	p, err := process.NewProcess(int32(pid))
	if err != nil {
		return nil, err
	}
	return parseProcess(p)
}

func ListProcesses() ([]Process, error) {
	processes, err := process.Processes()
	if err != nil {
		return nil, err
	}
	var result []Process
	for _, p := range processes {
		process, err := parseProcess(p)
		if err != nil {
			return nil, err
		}
		result = append(result, *process)
	}
	return result, nil
}

func parseProcess(p *process.Process) (*Process, error) {
	ppid, err := p.Ppid()
	if err != nil {
		return nil, err
	}
	var startTimePtr *time.Time
	startTimeMs, err := p.CreateTime()
	if err == nil {
		startTime := time.Unix(0, startTimeMs*int64(time.Millisecond))
		startTimePtr = &startTime
	}
	process := &Process{
		PID:  int(p.Pid),
		PPID: int(ppid),
	}
	process.StartTime = startTimePtr
	process.CommandLine, _ = p.Cmdline()
	process.Executable, _ = p.Exe()
	return process, nil
}
