package main

import (
	"github.com/jvikstedt/sneaky/engine"
)

type WaitRoom struct {
}

func NewWaitRoom(engine.RoomID, engine.CRHandler) engine.Room {
	return &WaitRoom{}
}
