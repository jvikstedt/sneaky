package engine

import (
	"io"
	"log"
	"net"
	"sync"
)

type connHandler interface {
	HandleConn(io.ReadWriteCloser) error
}

type TCPServer struct {
	logger   *log.Logger
	listener net.Listener
	ch       connHandler
}

func NewTCPServer(logger *log.Logger, ch connHandler) *TCPServer {
	return &TCPServer{
		logger: logger,
		ch:     ch,
	}
}

// Listen creates a tcp server and starts listening for connections on specified address
// When server is closed, this will block until all connections are done
func (s *TCPServer) Listen(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer l.Close()
	s.logger.Printf("Server running on %s, listening...", addr)

	s.listener = l
	var wg sync.WaitGroup
	for {
		conn, err := l.Accept()
		if err != nil {
			break
		}

		go func() {
			defer wg.Done()
			wg.Add(1)
			s.ch.HandleConn(conn)
		}()
	}

	s.logger.Println("Closing TCPServer...")
	wg.Wait()
	return nil
}

// Close closes server
func (s *TCPServer) Close() error {
	return s.listener.Close()
}
