package agent

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	DefaultAgentWorkdir = ".data"
)

var (
	DefaultAgentIdentityFilePath = filepath.Join(DefaultAgentWorkdir, "agent_identity.json")
)

const (
	ChannelSize = 10000
)

const (
	DefaultHeartbeatIntervalSeconds = 1
)

type Agent struct {
	Id                 string    `json:"id"`
	EphemeralId        string    `json:"ephemeral_id"`
	Time               time.Time `json:"time"`
	PID                int       `json:"pid,omitempty"`
	PPID               int       `json:"ppid,omitempty"`
	agentStatusChannel chan AgentStatus
}

func NewAgent() (*Agent, error) {
	err := os.MkdirAll(DefaultAgentWorkdir, 0700)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create .data directory")
	}
	agentIdentity, err := loadOrCreateAgentIdentity(DefaultAgentIdentityFilePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load or create agent identity")
	}
	agent := &Agent{
		Id:                 agentIdentity.Id,
		EphemeralId:        uuid.New().String(),
		Time:               time.Now(),
		PID:                os.Getpid(),
		PPID:               os.Getppid(),
		agentStatusChannel: make(chan AgentStatus, ChannelSize),
	}
	return agent, nil
}

func (a Agent) Run(ctx context.Context) error {
	var wg sync.WaitGroup

	wg.Add(1)

	go a.runEventLoop(ctx, &wg)

	wg.Wait()

	<-ctx.Done()

	return nil
}

func (a Agent) runEventLoop(ctx context.Context, wg *sync.WaitGroup) error {
	defer wg.Done()

	ticker := time.NewTicker(time.Duration(DefaultHeartbeatIntervalSeconds) * time.Second)

	for {
		select {
		case <-ctx.Done():
			log.Info("Exiting event loop")
			return nil
		case <-ticker.C:
			status := NewAgentStatus(AgentStatusTypeRunning)
			a.agentStatusChannel <- status
		case status := <-a.agentStatusChannel:
			log.Infof("Agent is %s (id=%s, ephemeral_id=%s, pid=%d, ppid=%d)", status.Status, a.Id, a.EphemeralId, a.PID, a.PPID)
		}
	}
}
