package main

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/whitfieldsdad/simplec2/pkg/shared"
)

func main() {
	ctx := context.Background()
	src := shared.NewProcessEventSource()
	ch, err := src.Run(ctx)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()

		for e := range ch {
			b, _ := json.Marshal(e)
			println(string(b))
		}
	}()

	wg.Wait()
}
