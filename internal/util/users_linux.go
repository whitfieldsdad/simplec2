package util

import (
	"context"
)

func listUsernames(ctx context.Context) ([]string, error) {
	command := "cut -d: -f1 /etc/passwd"
	return runCommandGetOutputSplitLines(ctx, command, "\n")
}

func listUserGroupNames(ctx context.Context) ([]string, error) {
	command := "cut -d: -f1 /etc/group"
	return runCommandGetOutputSplitLines(ctx, command, "\n")
}
