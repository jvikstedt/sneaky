package main

import (
	"encoding/json"

	"github.com/jvikstedt/sneaky/engine"
	"github.com/jvikstedt/sneaky/engine/packet"
)

type lobby struct {
	*engine.BaseRoom
}

func (l *lobby) New(id engine.RoomID, hub *engine.Hub) engine.Room {
	return &lobby{
		BaseRoom: engine.NewBaseRoom(id, hub),
	}
}

func (l *lobby) Name() engine.RoomName {
	return "lobby"
}
func (l *lobby) Initialize() {
}

type ChangeNameMessage struct {
	Username string
}

type SendMessage struct {
	Message string
}

func (l *lobby) HandleRequest(req engine.Request) error {
	switch req.CMD {
	case "CHANGE_NAME":
		msg := ChangeNameMessage{}
		err := json.Unmarshal(req.Payload, &msg)
		if err != nil {
			return err
		}
		req.Client.SaveValue("username", msg.Username)
	case "SEND_MESSAGE":
		msg := SendMessage{}
		err := json.Unmarshal(req.Payload, &msg)
		if err != nil {
			return err
		}

		username := "Anonymous"
		if i, err := req.Client.LoadValue("username"); err == nil {
			username = i.(string)
		}

		packet := packet.NewJsonMessage("NEW_MESSAGE", struct {
			Username string
			Message  string
		}{
			Username: username,
			Message:  msg.Message,
		}, nil)

		l.Broadcast(packet)
	}

	return nil
}
