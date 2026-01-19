package agent

import (
	"time"

	"github.com/google/uuid"
)

type AgentStatusType string

const (
	AgentStatusTypeRunning AgentStatusType = "running"
)

type AgentStatus struct {
	Id     string          `json:"id"`
	Time   time.Time       `json:"time"`
	Status AgentStatusType `json:"status"`
}

func NewAgentStatus(status AgentStatusType) AgentStatus {
	return AgentStatus{
		Id:     uuid.New().String(),
		Time:   time.Now(),
		Status: status,
	}
}
