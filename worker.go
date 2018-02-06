package main

import (
	"log"
	"image"
)

type RenderRequest struct {
	X, Y, Z int
	OutChan chan image.Image
}

func NewWorker(id int, workerQueue chan chan RenderRequest) Worker {
	worker := Worker{
		ID:          id,
		Work:        make(chan RenderRequest),
		WorkerQueue: workerQueue,
		QuitChan:    make(chan bool),
	}
	return worker
}

type Worker struct {
	ID          int
	Work        chan RenderRequest
	WorkerQueue chan chan RenderRequest
	QuitChan    chan bool
}

func (w *Worker) Start() {
	go func() {
		for {
			w.WorkerQueue <- w.Work

			select {
				case work := <-w.Work:
					log.Printf("worker%d: Received render request, X %d, Y, %d, Z %d\n", w.ID, work.X, work.Y, work.Z)
					work.OutChan <- image.NewNRGBA(image.Rect(0,0,1,1))

				case <-w.QuitChan:
					log.Printf("worker%d: stopping\n", w.ID)
					return
			}
		}
	}()
}

func (w *Worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}