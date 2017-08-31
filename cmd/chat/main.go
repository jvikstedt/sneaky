package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"

	"github.com/jvikstedt/sneaky/engine"
	"github.com/jvikstedt/sneaky/engine/packet"
)

func main() {
	cmd := os.Args[1]
	port := os.Getenv("CHAT_PORT")
	if port == "" {
		port = "8080"
	}

	switch cmd {
	case "client":
		runClient(port)
	default:
		runServer(port)
	}
}

func runServer(port string) {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	hub := engine.NewHub(logger, ":"+port)

	hub.RegisterRoom(&lobby{})

	err := hub.Start(engine.ServerTCP)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// CLIENT

type NewMessage struct {
	Username string
	Message  string
}

func runClient(port string) {
	var username string
	fmt.Print("Type username: ")
	fmt.Scanf("%s", &username)

	conn, err := net.Dial("tcp", ":"+port)
	if err != nil {
		panic(err)
	}

	msgCode := packet.NewMsgCode(conn, conn)
	reqCh := make(chan engine.Request)

	logger := log.New(ioutil.Discard, "", log.LstdFlags)
	client := engine.NewClient(1, logger, conn, msgCode, reqCh)

	go func() {
		for {
			select {
			case req := <-reqCh:
				if req.CMD == "NEW_MESSAGE" {
					var msg NewMessage
					json.Unmarshal(req.Payload, &msg)
					fmt.Printf("%s: %s\n", msg.Username, msg.Message)
				}
			}
		}
	}()

	client.Write(packet.Message{CMD: "CHANGE_NAME", Payload: []byte(`{"username": "` + username + `"}`)})

	reader := bufio.NewReader(os.Stdin)
	go func() {
		for {
			message, _ := reader.ReadString('\n')
			message = strings.TrimSpace(message)

			client.Write(packet.Message{CMD: "SEND_MESSAGE", Payload: []byte(`{"message": "` + message + `"}`)})
		}
	}()

	client.Start()
}
