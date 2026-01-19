package agent

import "time"

type Task struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreateTime  time.Time `json:"create_time"`
}

type TaskInvocation struct {
	Id         string    `json:"id"`
	HostId     string    `json:"host_id"`
	AgentId    string    `json:"agent_id"`
	CreateTime time.Time `json:"create_time"`
}

type TaskInvocationResult struct {
	Id               string `json:"id"`
	TaskId           string `json:"task_id"`
	TaskInvocationId string `json:"task_invocation_id"`
}
