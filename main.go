package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"strconv"
	"runtime"
	"os"
	"io/ioutil"
)

func handleTile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Make tile
	//   No need to check for invalid ints as this in done by regex
	z, _ := strconv.Atoi(vars["z"])
	x, _ := strconv.Atoi(vars["x"])
	y, _ := strconv.Atoi(vars["y"])
	tile := NewTileWithXY(x, y, z)

	var img []byte
		out := make(chan []byte)
		request := RenderRequest{tile, out}

		// Request render
		log.Printf("Reqeust render of X %d Y %d Z %d\n", x, y, z)
		RenderQueue <- request

	f, err := os.Create("image.png")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
		// Wait for render
		img = <-out

	err = ioutil.WriteFile("image.png", img, 0644)
	if err != nil {
		log.Fatal(err)
	}

	// Send image to client
	w.Header().Add("Content-Type", "image/png")
	w.Write(img)
}

func main() {
	log.Println("Starting the dispatcher")
	StartDispatcher(runtime.NumCPU())

	r := mux.NewRouter()
	r.Methods("GET").Path("/osm_tiles/{z:[0-9]+}/{x:[0-9]+}/{y:[0-9]+}.png").HandlerFunc(handleTile)

	http.ListenAndServe(":8001", r)
}
