package agent

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/whitfieldsdad/simplec2/internal/util"
)

const (
	ChannelSize       = 1000
	HeartbeatInterval = 1 * time.Second
)

var (
	AgentDataDir       = filepath.Join("data", "agents")
	AgentStateFilePath = filepath.Join("agent-state.json")
)

type AgentState struct {
	HostId     string          `json:"host_id"`
	Identities []util.Identity `json:"identities,omitempty"`
}

func NewAgentState(ctx context.Context) (*AgentState, error) {
	hostId := uuid.New().String()
	identities, err := util.ListIdentities(ctx)
	if err != nil {
		return nil, err
	}
	state := &AgentState{
		HostId:     hostId,
		Identities: identities,
	}
	return state, nil
}

type Agent struct {
	Id                 string `json:"id"`
	agentStatusChannel chan AgentStatus
}

func NewAgent() (*Agent, error) {
	agent := &Agent{
		Id:                 uuid.New().String(),
		agentStatusChannel: make(chan AgentStatus, ChannelSize),
	}
	return agent, nil
}

func GetAgentState(ctx context.Context) (*AgentState, error) {
	b, err := os.ReadFile(AgentStateFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			state, err := NewAgentState(ctx)
			if err != nil {
				return nil, err
			}
			b, err := json.MarshalIndent(state, "", "  ")
			if err != nil {
				return nil, err
			}
			err = os.WriteFile(AgentStateFilePath, b, 0o644)
			if err != nil {
				return nil, err
			}
			return state, nil
		}
	}

	var state AgentState
	err = json.Unmarshal(b, &state)
	if err != nil {
		return nil, err
	}
	return &state, nil
}

func (a Agent) Run(ctx context.Context) error {
	var (
		wg  sync.WaitGroup
		err error
	)
	pid := os.Getpid()
	ppid := os.Getppid()

	log.Printf("Starting agent")

	// Read agent state
	agentState, err := GetAgentState(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get agent state")
	}
	hostId := agentState.HostId

	a.agentStatusChannel <- NewAgentStatus(hostId, a.Id, AgentStatusTypeStarting)

	agentDataDir := filepath.Join(AgentDataDir, hostId, a.Id)
	agentStatusesFilePath := filepath.Join(agentDataDir, "agent-statuses.jsonl")

	// Create `data/agents/<host_id>/<agent_id>/` directory if it doesn't exist.
	err = os.MkdirAll(agentDataDir, 0o755)
	if err != nil {
		return err
	}

	// Continuously write agent statuses to `data/agents/<host_id>/<agent_id>/agent-statuses.jsonl`.
	go func() {
		log.Printf("Writing agent statuses to %s", agentStatusesFilePath)
		f, err := os.OpenFile(agentStatusesFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
		if err != nil {
			log.Fatalf("failed to open agent statuses file: %v", err)
		}
		defer f.Close()

		for status := range a.agentStatusChannel {
			log.Printf("Received agent status update (agent_id=%s, pid=%d, ppid=%d, status=%s)", status.AgentId, pid, ppid, status.Status)
			b, err := json.Marshal(status)
			if err != nil {
				log.Fatalf("failed to marshal agent status: %v", err)
			}
			_, err = f.Write(append(b, '\n'))
		}
	}()

	wg.Add(1)

	go a.runAgentHeartbeat(ctx, hostId, &wg)

	log.Printf("Agent is running")
	a.agentStatusChannel <- NewAgentStatus(hostId, a.Id, AgentStatusTypeRunning)

	wg.Wait()

	log.Printf("Stopping agent")
	a.agentStatusChannel <- NewAgentStatus(hostId, a.Id, AgentStatusTypeStopping)

	close(a.agentStatusChannel)

	log.Printf("Agent stopped")
	return nil
}

func (a Agent) runAgentHeartbeat(ctx context.Context, hostId string, wg *sync.WaitGroup) {
	defer wg.Done()

	a.agentStatusChannel <- NewAgentStatus(hostId, a.Id, AgentStatusTypeHeartbeat)

	ticker := time.NewTicker(HeartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			a.agentStatusChannel <- NewAgentStatus(hostId, a.Id, AgentStatusTypeHeartbeat)
		}
	}
}
