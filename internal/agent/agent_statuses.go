package agent

import (
	"time"

	"github.com/google/uuid"
)

type AgentStatusType string

const (
	AgentStatusTypeStarting  AgentStatusType = "starting"
	AgentStatusTypeRunning   AgentStatusType = "running"
	AgentStatusTypeHeartbeat AgentStatusType = "heartbeat"
	AgentStatusTypeStopping  AgentStatusType = "stopping"
)

type AgentStatus struct {
	Id      string          `json:"id"`
	HostId  string          `json:"host_id"`
	AgentId string          `json:"agent_id"`
	Time    time.Time       `json:"time"`
	Status  AgentStatusType `json:"status"`
}

func NewAgentStatus(hostId, agentId string, status AgentStatusType) AgentStatus {
	return AgentStatus{
		Id:      uuid.New().String(),
		Time:    time.Now(),
		HostId:  hostId,
		AgentId: agentId,
		Status:  status,
	}
}
