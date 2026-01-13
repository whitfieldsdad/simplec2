package util

import (
	"context"
	"os"
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

	dbusUUID, err := getDBUSUUID(ctx)
	if err == nil && dbusUUID != "" {
		identities = append(identities, Identity{
			Type:  IdentityTypeDBUSUUID,
			Value: dbusUUID,
		})
	}

	return identities, nil
}

func getSMBIOSUUID(ctx context.Context) (string, error) {
	command := "dmidecode -s system-uuid"
	output, err := runCommandGetOutput(ctx, command)
	if err != nil {
		return "", err
	}
	uuid := strings.ToLower(strings.Trim(output, "\n "))
	return uuid, nil
}

func getDBUSUUID(ctx context.Context) (string, error) {
	paths := []string{
		"/var/lib/dbus/machine-id",
		"/etc/machine-id",
	}
	for _, path := range paths {
		b, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		uuid := strings.Split(string(b), "\n")[0]
		return uuid, nil
	}
	return "", nil
}
