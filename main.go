package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"strconv"
	"runtime"
	"os"
	"path/filepath"
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

	// Check cache
	var img []byte
	file, err := os.Open(filepath.Join("./cache", tile.URL()))
	if err != nil {
		// Make render request
		out := make(chan []byte)
		request := RenderRequest{tile, out}

		// Request render
		log.Printf("Reqeust render of X %d Y %d Z %d\n", x, y, z)
		RenderQueue <- request

		// Wait for render
		img = <-out

		// Save to cache
		err := os.MkdirAll(filepath.Join("./cache", tile.URLBase()), 0755)
		if err != nil {
			log.Printf("Error creating tile folder: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		file, err = os.OpenFile(filepath.Join("./cache", tile.URL()), os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Printf("Error creating tile file: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err = file.Write(img)
		if err != nil {
			log.Printf("Error saving tile: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		// Load from cache
		var err error
		img, err = ioutil.ReadAll(file)
		if err != nil {
			log.Printf("Error reading tile: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	// Send image to client
	w.Header().Add("Content-Type", "image/png")
	w.Write(img)
}

func main() {
	log.Println("Setting up cache")
	os.Mkdir("./cache", 0755)

	log.Println("Starting the dispatcher")
	StartDispatcher(runtime.NumCPU())

	r := mux.NewRouter()
	r.Methods("GET").Path("/osm_tiles/{z:[0-9]+}/{x:[0-9]+}/{y:[0-9]+}.png").HandlerFunc(handleTile)

	http.ListenAndServe(":8001", r)
}
