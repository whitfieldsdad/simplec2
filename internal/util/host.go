package util

type Host struct {
	Id                string             `json:"id"`
	Identities        []Identity         `json:"identities,omitempty"`
	OperatingSystem   *OperatingSystem   `json:"operating_system,omitempty"`
	NetworkInterfaces []NetworkInterface `json:"network_interfaces,omitempty"`
}

func (h Host) GetArtifactType() ArtifactType {
	return ArtifactTypeHost
}
