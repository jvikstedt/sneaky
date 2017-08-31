package main

import (
	"fmt"

	"github.com/jvikstedt/sneaky/engine"
)

type waitRoom struct {
	*engine.BaseRoom
}

func (w *waitRoom) New(id engine.RoomID, hub *engine.Hub) engine.Room {
	return &waitRoom{
		BaseRoom: engine.NewBaseRoom(id, hub),
	}
}

func (w *waitRoom) Name() engine.RoomName {
	return "wait"
}

func (w *waitRoom) Initialize() {
	w.AddAction("PING", w.ping)
}

func (w *waitRoom) ping(engine.Request) error {
	fmt.Println("ping called")
	return nil
}
