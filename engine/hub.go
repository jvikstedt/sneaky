package engine

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/jvikstedt/sneaky/engine/packet"
)

type server interface {
	Listen(string) error
	Close() error
}

type ServerType int

const (
	ServerTCP ServerType = iota
	ServerWS
)

type Hub struct {
	logger *log.Logger
	server
	addr string

	rLock           sync.RWMutex
	registeredRooms map[RoomName]Room
	rooms           map[RoomID]Room

	cLock   sync.RWMutex
	clients map[ClientID]*Client

	reqCh chan Request
}

func NewHub(logger *log.Logger, addr string) *Hub {
	return &Hub{
		logger:          logger,
		addr:            addr,
		registeredRooms: make(map[RoomName]Room),
		rooms:           make(map[RoomID]Room),
		clients:         make(map[ClientID]*Client),
		reqCh:           make(chan Request, 20),
	}
}

func (h *Hub) Start(t ServerType) error {
	switch t {
	case ServerTCP:
		h.server = NewTCPServer(h.logger, h)
	case ServerWS:
		h.server = NewWsServer(h.logger, h)
	default:
		return fmt.Errorf("Unknown ServerType %d", t)
	}

	onSignal(func(sig os.Signal) {
		if sig == os.Interrupt {
			h.Shutdown()
		}
	})

	go func() {
		for {
			select {
			case req := <-h.reqCh:
				req.Client.Room().HandleRequest(req)
			}
		}
	}()

	return h.server.Listen(h.addr)
}

func (h *Hub) Shutdown() error {
	defer h.closeClients()
	return h.server.Close()
}

// RegisterRoom adds room to list of available rooms
// First registered room will also be used as a join room
func (h *Hub) RegisterRoom(r Room) {
	h.registeredRooms[r.Name()] = r
	if len(h.rooms) == 0 {
		h.NewRoom(r.Name())
	}
}

func (h *Hub) NewRoom(name RoomName) {
	h.rLock.Lock()
	defer h.rLock.Unlock()
	id := RoomID(nextID())
	h.rooms[id] = h.registeredRooms[name].New(id, h)
	h.rooms[id].Initialize()
}

func (h *Hub) getJoinRoom() Room {
	h.rLock.Lock()
	defer h.rLock.Unlock()
	return h.rooms[1]
}

func (h *Hub) HandleConn(conn io.ReadWriteCloser) error {
	client := NewClient(ClientID(nextID()), h.logger, conn, packet.NewMsgCode(conn, conn), h.reqCh)
	defer h.clientCleanup(client)

	err := h.MigrateClientToRoom(client, h.getJoinRoom())
	if err != nil {
		return err
	}

	h.addClient(client)

	return client.Start()
}

func (h *Hub) addClient(c *Client) {
	h.cLock.Lock()
	h.clients[c.ID()] = c
	h.cLock.Unlock()
}

func (h *Hub) removeClient(c *Client) {
	h.cLock.Lock()
	delete(h.clients, c.ID())
	h.cLock.Unlock()
}

func (h *Hub) closeClients() {
	h.cLock.Lock()
	for _, c := range h.clients {
		c.Close()
	}
	h.cLock.Unlock()
}

func (h *Hub) MigrateClientToRoom(c *Client, r Room) error {
	if c.Room() != nil {
		c.Room().RemoveClient(c)
	}

	err := r.AddClient(c)
	if err != nil {
		return err
	}

	c.setRoom(r)

	return nil
}

func (h *Hub) clientCleanup(c *Client) {
	c.Close()

	h.removeClient(c)
	c.Room().RemoveClient(c)
}

func onSignal(f func(os.Signal)) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		f(sig)
	}()
}
