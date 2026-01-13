package util

import (
	"time"

	"github.com/google/uuid"
)

type ArtifactType string

const (
	ArtifactTypeHost             ArtifactType = "host"
	ArtifactTypeFile             ArtifactType = "file"
	ArtifactTypeProcess          ArtifactType = "process"
	ArtifactTypeUser             ArtifactType = "user"
	ArtifactTypeUserGroup        ArtifactType = "user-group"
	ArtifactTypeOperatingSystem  ArtifactType = "operating-system"
	ArtifactTypeNetworkInterface ArtifactType = "network-interface"
)

type Artifact struct {
	Id   string            `json:"id"`
	Time time.Time         `json:"time"`
	Type ArtifactType      `json:"type"`
	Data ArtifactInterface `json:"data"`
}

func NewArtifact(data ArtifactInterface) Artifact {
	return Artifact{
		Id:   uuid.New().String(),
		Time: time.Now(),
		Type: data.GetArtifactType(),
		Data: data,
	}
}

type ArtifactInterface interface {
	GetArtifactType() ArtifactType
}

type ArtifactBundle struct {
	Id        string     `json:"id"`
	Time      time.Time  `json:"time"`
	Artifacts []Artifact `json:"artifacts"`
}

func NewArtifactBundle() ArtifactBundle {
	return ArtifactBundle{
		Id:   uuid.New().String(),
		Time: time.Now(),
	}
}
