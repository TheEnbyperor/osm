package main

import (
	"log"
	"time"
	"image"
)

func main() {
	log.Println("Starting the dispatcher")
	StartDispatcher(3)

	time.Sleep(time.Second * 2)

	out := make(chan image.Image)
	request := RenderRequest{2011, 1362, 12, out}

	RenderQueue <- request

	log.Println("Waiting for image back")
	img := <- out
	log.Println(img)
}
