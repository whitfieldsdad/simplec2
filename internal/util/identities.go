package util

import "context"

type IdentityType string

const (
	IdentityTypeSimpleC2AgentUUID IdentityType = "simplec2-agent-uuid"
	IdentityTypeSMBIOSUUID        IdentityType = "smbios-uuid"
	IdentityTypeDBUSUUID          IdentityType = "dbus-uuid"
)

var (
	IdentityPriorityOrder = []IdentityType{
		IdentityTypeSMBIOSUUID,
		IdentityTypeDBUSUUID,
	}
)

type Identity struct {
	Type  IdentityType `json:"type"`
	Value string       `json:"value"`
}

func ListIdentities(ctx context.Context) ([]Identity, error) {
	return listIdentities(ctx)
}
