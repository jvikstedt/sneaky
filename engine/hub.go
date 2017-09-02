package engine

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
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
}

func NewHub(logger *log.Logger, addr string) *Hub {
	return &Hub{
		logger: logger,
		addr:   addr,
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

	//go func() {
	//	for {
	//		select {
	//		case req := <-h.reqCh:
	//			req.Client.Room().HandleRequest(req)
	//		}
	//	}
	//}()

	return h.server.Listen(h.addr)
}

func (h *Hub) Shutdown() error {
	return h.server.Close()
}

func (h *Hub) HandleConn(conn io.ReadWriteCloser) error {
	//client := NewClient(ClientID(nextID()), h.logger, conn, packet.NewMsgCode(conn, conn), h.reqCh)
	//defer h.clientCleanup(client)

	////err := h.MigrateClientToRoom(client, h.getJoinRoom())
	////if err != nil {
	////	return err
	////}

	//h.addClient(client)

	//return client.Start()
	return nil
}

func onSignal(f func(os.Signal)) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		f(sig)
	}()
}
