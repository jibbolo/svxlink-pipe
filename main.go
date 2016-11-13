package main

import (
	"net/http"
	"os"
)

var pipe *Pipe

const maxSize = 30
const maxRowBytes = 500

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3334"
	}
	pipe = NewPipe()
	go pipe.Scan(os.Stdin)
	http.ListenAndServe(":"+port, pipe.NewRouter())
}
