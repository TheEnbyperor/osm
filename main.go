package main

import (
	"log"
	"os"
	"io/ioutil"
)

func main() {
	log.Println("Starting the dispatcher")
	StartDispatcher(3)

	out := make(chan []byte)
	tile := NewTileWithXY(8047, 5449, 14)
	request := RenderRequest{tile, out}

	RenderQueue <- request

	log.Println("Waiting for image back")
	img := <- out
	log.Println("Got image saving")

	f, err := os.Create("image.png")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	err = ioutil.WriteFile("image.png", img, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
