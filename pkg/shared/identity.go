package shared

import (
	"github.com/denisbrodbeck/machineid"
	"github.com/google/uuid"
)

func GetHostId() string {
	id, err := machineid.ProtectedID(ProductId)
	if err != nil {
		panic(err)
	}
	return id
}

func GetAgentId() string {
	return uuid.New().String()
}
