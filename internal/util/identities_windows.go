package util

import (
	"context"
	"strings"
)

func listIdentities(ctx context.Context) ([]Identity, error) {
	identities := []Identity{}

	smbiosUUID, err := getSMBIOSUUID(ctx)
	if err == nil && smbiosUUID != "" {
		identities = append(identities, Identity{
			Type:  IdentityTypeSMBIOSUUID,
			Value: smbiosUUID,
		})
	}

	return identities, nil
}

func getSMBIOSUUID(ctx context.Context) (string, error) {
	command := "(Get-CimInstance -Class Win32_ComputerSystemProduct).UUID"
	output, err := runCommandGetOutput(ctx, "powershell -Command \""+command+"\"")
	if err != nil {
		return "", err
	}
	uuid := strings.ToLower(strings.Trim(output, "\n "))
	return uuid, nil
}
