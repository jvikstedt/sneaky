package engine

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type WsServer struct {
	logger *log.Logger
	ch     connHandler
	server *http.Server
	wg     sync.WaitGroup
}

func NewWsServer(logger *log.Logger, ch connHandler) *WsServer {
	return &WsServer{
		logger: logger,
		ch:     ch,
	}
}

func (s *WsServer) Listen(addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", s.ws)

	s.server = &http.Server{Addr: addr, Handler: mux}
	err := s.server.ListenAndServe()
	s.wg.Wait()
	return err
}

func (s *WsServer) Close() error {
	return s.server.Shutdown(nil)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *WsServer) ws(w http.ResponseWriter, r *http.Request) {
	s.wg.Add(1)
	defer s.wg.Done()

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.Println(err)
		return
	}
	defer conn.Close()

	wsConn := &WSConnection{
		conn: conn,
	}

	s.ch.HandleConn(wsConn)
}

type WSConnection struct {
	conn *websocket.Conn
}

func (c *WSConnection) Close() error {
	return c.conn.Close()
}

func (c *WSConnection) Read(p []byte) (int, error) {
	_, d, err := c.conn.ReadMessage()
	copy(p, d)
	return len(d), err
}

func (c *WSConnection) Write(p []byte) (int, error) {
	err := c.conn.WriteMessage(websocket.TextMessage, p)
	return len(p), err
}
