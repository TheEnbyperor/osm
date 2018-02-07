package main

import (
	"log"
	"image"
	"github.com/fawick/go-mapnik/mapnik"
	"image/png"
	"bytes"
)

type RenderRequest struct {
	Tile *Tile
	OutChan chan []byte
}

func NewWorker(id int, workerQueue chan chan RenderRequest) Worker {
	m := mapnik.NewMap(256, 256)
	err := m.Load("OSMBright/OSMBright.xml")
	if err != nil {
		log.Println(err)
	}
	worker := Worker{
		ID:          id,
		Work:        make(chan RenderRequest),
		WorkerQueue: workerQueue,
		QuitChan:    make(chan bool),
		m:           m,
	}
	return worker
}

type Worker struct {
	ID          int
	Work        chan RenderRequest
	WorkerQueue chan chan RenderRequest
	QuitChan    chan bool
	m           *mapnik.Map
}

func (w *Worker) Start() {
	go func() {
		for {
			w.WorkerQueue <- w.Work

			select {
				case work := <-w.Work:
					log.Printf("worker%d: Received render request, X %d, Y, %d, Z %d\n", w.ID,
						work.Tile.X, work.Tile.Y, work.Tile.Z)

					img, err := w.renderTile(work.Tile)
					if err != nil {
						log.Println("Error rendering")
						buf := new(bytes.Buffer)
						png.Encode(buf, image.NewNRGBA(image.Rect(0,0,1,1)))
						work.OutChan <- buf.Bytes()
					} else {
						work.OutChan <- img
					}


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

func (w *Worker) renderTile(t *Tile) ([]byte, error) {
	next := t.NextTile()
	p := w.m.Projection()
	ll := p.Forward(mapnik.Coord{t.Long, t.Lat})
	ur := p.Forward(mapnik.Coord{next.Long, next.Lat})
	log.Println(t, next)
	w.m.ZoomToMinMax(ll.X, ll.Y, ur.X, ur.Y)

	img, err := w.m.RenderToMemoryPng()
	if err != nil {
		return nil, err
	} else {
		return img, nil
	}
}