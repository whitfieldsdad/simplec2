package util

import (
	"context"
)

func listUsernames(ctx context.Context) ([]string, error) {
	command := "dscl . list /Users"
	return runCommandGetOutputSplitLines(ctx, command, "\n")
}

func listUserGroupNames(ctx context.Context) ([]string, error) {
	command := "dscl . list /Groups"
	return runCommandGetOutputSplitLines(ctx, command, "\n")
}
