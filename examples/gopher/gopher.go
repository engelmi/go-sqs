package main

import (
	"sync"
	"time"

	"github.com/engelmi/go-sqs/examples"
)

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go examples.NewGopher("gopher_queue", "ConsumingGopher-1", 100*time.Millisecond).Consume()
	go examples.NewGopher("gopher_queue", "PushingGopher-1", 1*time.Second).StartPushingMessages()
	go examples.NewGopher("gopher_queue", "PushingGopher-2", 2*time.Second).StartPushingMessages()
	go examples.NewGopher("gopher_queue", "PushingGopher-3", 1*time.Second).StartPushingMessages()
	go examples.NewGopher("gopher_queue", "PushingGopher-4", 4*time.Second).StartPushingMessages()

	wg.Wait()
}
