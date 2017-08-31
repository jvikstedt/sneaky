package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jvikstedt/sneaky/engine"
)

func main() {
	port := os.Getenv("SNEAKY_PORT")
	if port == "" {
		port = "8080"
	}

	logger := log.New(os.Stdout, "", log.LstdFlags)
	hub := engine.NewHub(logger, ":"+port)

	hub.RegisterRoom(&waitRoom{})

	err := hub.Start(engine.ServerWS)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
