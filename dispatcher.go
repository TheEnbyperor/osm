package main

import "log"

var RenderQueue = make(chan RenderRequest, 100)

var WorkerQueue chan chan RenderRequest

func StartDispatcher(nworkers int) {
	WorkerQueue = make(chan chan RenderRequest, nworkers)

	for i := 0; i<nworkers; i++ {
		log.Println("Starting worker", i+1)
		worker := NewWorker(i+1, WorkerQueue)
		worker.Start()
	}

	go func() {
		for {
			select {
			case work := <-RenderQueue:
				log.Println("Received render requeust")
				go func() {
					worker := <-WorkerQueue

					log.Println("Dispatching render request")
					worker <- work
				}()
			}
		}
	}()
}