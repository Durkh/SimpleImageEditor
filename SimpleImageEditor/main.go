package main

import (
	"SimpleImageEditor/backend"
	"SimpleImageEditor/common"
	"SimpleImageEditor/frontend/cli"
	"SimpleImageEditor/queue"
	"log"
	"sync"
)

func main() {

	var (
		data          = make(chan common.Info, 1)
		q, closeQueue = queue.New(32)
		wg            sync.WaitGroup
		done          = make(chan bool)
	)

	wg.Add(2)
	go cli.Run(q.Sender(), data, done, &wg)

	select {
	case <-done:
		log.Println("done")
		closeQueue()
	}

	go backend.Run(q.Receiver(), data, &wg)
	wg.Wait()

}
