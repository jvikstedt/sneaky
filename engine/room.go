package engine

import (
	"sync"

	"github.com/jvikstedt/sneaky/engine/packet"
)

type RoomName string
type RoomID int

type Room interface {
	New(RoomID, *Hub) Room
	Name() RoomName
	ID() RoomID
	HandleRequest(Request) error
	AddClient(*Client) error
	RemoveClient(*Client)
	Initialize()
}

type actionHandler func(Request) error

type BaseRoom struct {
	id  RoomID
	hub *Hub

	muClients sync.RWMutex
	clients   map[ClientID]*Client
	actions   map[packet.Command]actionHandler
}

func NewBaseRoom(id RoomID, hub *Hub) *BaseRoom {
	return &BaseRoom{
		id:      id,
		hub:     hub,
		clients: make(map[ClientID]*Client),
		actions: make(map[packet.Command]actionHandler),
	}
}

func (b *BaseRoom) ID() RoomID {
	return b.id
}

func (b *BaseRoom) Hub() *Hub {
	return b.hub
}

func (b *BaseRoom) AddClient(c *Client) error {
	b.muClients.Lock()
	b.clients[c.ID()] = c
	b.muClients.Unlock()

	return nil
}

func (b *BaseRoom) RemoveClient(c *Client) {
	b.muClients.Lock()
	delete(b.clients, c.ID())
	b.muClients.Unlock()
}

func (b *BaseRoom) Clients(f func(map[ClientID]*Client)) {
	b.muClients.RLock()
	f(b.clients)
	b.muClients.RUnlock()
}

func (b *BaseRoom) Broadcast(msg packet.Message) {
	b.Clients(func(clients map[ClientID]*Client) {
		for _, v := range clients {
			v.Write(msg)
		}
	})
}

func (b *BaseRoom) Action(cmd packet.Command) actionHandler {
	return b.actions[cmd]
}

func (b *BaseRoom) AddAction(cmd packet.Command, action actionHandler) {
	b.actions[cmd] = action
}

func (b *BaseRoom) HandleRequest(req Request) error {
	return b.Action(req.CMD)(req)
}
