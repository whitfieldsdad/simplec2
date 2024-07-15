package shared

import (
	"context"
	"time"

	"github.com/shirou/gopsutil/v3/process"
)

type ProcessIdentity struct {
	Pid  int32 `json:"pid"`
	Ppid int32 `json:"ppid"`
}

func (i ProcessIdentity) GetUUID() string {
	m := map[string]interface{}{
		"pid":  i.Pid,
		"ppid": i.Ppid,
	}
	return CalculateUUIDv5FromMap(m)
}

type Process struct {
	StartTime   *time.Time `json:"start_time,omitempty"`
	Pid         int32      `json:"pid"`
	Ppid        int32      `json:"ppid"`
	Name        string     `json:"name,omitempty"`
	Username    string     `json:"user_id,omitempty"`
	Executable  *File      `json:"executable,omitempty"`
	CommandLine []string   `json:"command_line,omitempty"`
}

func (p Process) GetIdentity() ProcessIdentity {
	return ProcessIdentity{
		Pid:  p.Pid,
		Ppid: p.Ppid,
	}
}

func (p Process) GetObservableType() ObservableType {
	return ObservableTypeProcess
}

type ProcessQuery struct {
	Pids []int32 `json:"pids,omitempty"`
}

func (q ProcessQuery) Matches(p *process.Process) bool {
	if len(q.Pids) == 0 {
		return true
	}
	for _, pid := range q.Pids {
		if pid == p.Pid {
			return true
		}
	}
	return false
}

func ListPids(ctx context.Context) ([]int32, error) {
	processes, err := process.ProcessesWithContext(ctx)
	if err != nil {
		return nil, err
	}
	pids := make([]int32, 0, len(processes))
	for _, p := range processes {
		pids = append(pids, p.Pid)
	}
	return pids, nil
}

func ListProcesses(ctx context.Context, q *ProcessQuery) (chan Process, error) {
	ch := make(chan Process)
	go func() {
		defer close(ch)

		processes, err := process.ProcessesWithContext(ctx)
		if err != nil {
			return
		}
		for _, p := range processes {
			if q != nil && !q.Matches(p) {
				continue
			}
			process, err := getProcess(ctx, p)
			if err != nil {
				continue
			}
			ch <- process
		}
	}()
	return ch, nil
}

func getProcess(ctx context.Context, p *process.Process) (Process, error) {
	ppid, err := p.Ppid()
	if err != nil {
		return Process{}, err
	}
	var startTimePtr *time.Time
	startTimeMs, err := p.CreateTime()
	if err == nil {
		startTime := time.Unix(0, startTimeMs*int64(time.Millisecond))
		startTimePtr = &startTime
	}
	process := Process{
		Pid:  p.Pid,
		Ppid: ppid,
	}
	process.StartTime = startTimePtr
	process.Name, _ = p.NameWithContext(ctx)
	process.Username, _ = p.UsernameWithContext(ctx)
	process.CommandLine, _ = p.CmdlineSliceWithContext(ctx)

	exe, err := p.ExeWithContext(ctx)
	if err == nil && exe != "" {
		process.Executable, _ = GetFile(exe)
	}
	return process, nil
}
