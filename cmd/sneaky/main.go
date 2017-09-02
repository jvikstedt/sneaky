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

	cr := engine.NewCRHandler()
	cr.RegisterCreator("WaitRoom", NewWaitRoom, 4)

	logger := log.New(os.Stdout, "", log.LstdFlags)
	hub := engine.NewHub(logger, ":"+port)

	err := hub.Start(engine.ServerWS)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
