package main

import (
	"log"

	"pix-psp-simulator/src/containers"
)

func main() {
	cfg := containers.LoadConfig()
	server := containers.NewServer(cfg)
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
