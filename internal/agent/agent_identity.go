package agent

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

type AgentIdentity struct {
	Id   string    `json:"id"`   // Persistent, globally unique identifier for the agent
	Time time.Time `json:"time"` // When the agent identity was created
}

func (i *AgentIdentity) Read(r io.Reader) error {
	dec := json.NewDecoder(r)
	err := dec.Decode(i)
	if err != nil {
		return err
	}
	return nil
}

func (i *AgentIdentity) Write(w io.Writer) error {
	enc := json.NewEncoder(w)
	err := enc.Encode(i)
	if err != nil {
		return err
	}
	return nil
}

func newAgentIdentity() *AgentIdentity {
	return &AgentIdentity{
		Id:   uuid.New().String(),
		Time: time.Now(),
	}
}

func loadOrCreateAgentIdentity(path string) (*AgentIdentity, error) {
	var identity *AgentIdentity

	f, err := os.Open(path)
	if err == nil {
		defer f.Close()

		log.Infof("Reading agent identity file: %s", path)

		identity = &AgentIdentity{}
		err = identity.Read(f)
		if err != nil {
			return nil, err
		}

	} else {
		if errors.Is(err, os.ErrNotExist) {
			log.Infof("Generating agent identity file: %s", path)

			identity = newAgentIdentity()
			file, err := os.Create(path)
			if err != nil {
				return nil, err
			}
			defer file.Close()

			err = identity.Write(file)
			if err != nil {
				return nil, err
			}
		}
	}
	return identity, nil
}
