package main

import (
	"log"
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
	pipe = NewPipe(port, os.Stdin)
	log.Fatal(pipe.Run())
}
