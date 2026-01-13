package util

import (
	"context"
)

func listUsernames(ctx context.Context) ([]string, error) {
	command := "Get-LocalUser | Select-Object -ExpandProperty Name"
	return runCommandGetOutputSplitLines(ctx, command, "\n")
}

func listUserGroupNames(ctx context.Context) ([]string, error) {
	command := "Get-LocalGroup | Select-Object -ExpandProperty Name"
	return runCommandGetOutputSplitLines(ctx, command, "\n")
}
