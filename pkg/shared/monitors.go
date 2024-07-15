package shared

import (
	"context"
	"fmt"
	"slices"
	"sync"
	"time"
)

type EventSource interface {
	Run(context.Context) (chan Event, error)
}

type ProcessEventSource struct{}

func NewProcessEventSource() *ProcessEventSource {
	return &ProcessEventSource{}
}

func (m *ProcessEventSource) Run(ctx context.Context) (chan Event, error) {
	var oldPids []int32
	oldProcesses, err := ListProcesses(ctx, nil)
	if err != nil {
		return nil, err
	}
	for p := range oldProcesses {
		oldPids = append(oldPids, p.Pid)
	}

	ch := make(chan Event)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(10 * time.Millisecond):
				pids, err := ListPids(ctx)
				if err != nil {
					continue
				}
				var startedPids []int32
				for _, pid := range pids {
					if !slices.Contains(oldPids, pid) {
						startedPids = append(startedPids, pid)
					}
				}
				var stoppedPids []int32
				for _, oldPid := range oldPids {
					if !slices.Contains(pids, oldPid) {
						stoppedPids = append(stoppedPids, oldPid)
					}
				}
				if len(startedPids) > 0 {
					fmt.Printf("%d processes started\n", len(startedPids))
				}
				if len(stoppedPids) > 0 {
					fmt.Printf("%d processes stopped\n", len(stoppedPids))
				}
				oldPids = pids
			}
		}
	}()

	return ch, nil
}
