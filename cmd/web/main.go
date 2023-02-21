package main

import (
	"log"
	"net/http"
)

func main() {
	path := "web/"
	fs := http.FileServer(http.Dir(path))
	http.Handle("/", fs)

	log.Print("Listening on :3000...")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
