package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/whitfieldsdad/simplec2/internal/agent"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	agent, err := agent.NewAgent()
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}
	err = agent.Run(ctx)
	if err != nil {
		log.Fatalf("Agent failed: %v", err)
	}
}
