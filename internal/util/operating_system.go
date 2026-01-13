package util

import "runtime"

type OperatingSystem struct {
	Type string `json:"type"`
}

func (o OperatingSystem) GetArtifactType() ArtifactType {
	return ArtifactTypeOperatingSystem
}

func GetOperatingSystem() OperatingSystem {
	return OperatingSystem{
		Type: runtime.GOOS,
	}
}
