package util

import (
	"context"
	"log"
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
	command := "system_profiler SPHardwareDataType | awk '/UUID/ { print $3; }'"
	output, err := runCommandGetOutput(ctx, command)
	if err != nil {
		return "", err
	}
	log.Printf("SMBIOS UUID: %s", output)
	uuid := strings.ToLower(strings.Trim(output, "\n "))
	return uuid, nil
}
