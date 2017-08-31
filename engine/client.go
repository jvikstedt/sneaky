package engine

import (
	"errors"
	"io"
	"log"
	"sync"

	"github.com/jvikstedt/sneaky/engine/packet"
)

type ClientID int

// Client provides easy way to read/write Message's
// from/to MsgCoder
type Client struct {
	// id can be used to keep track which client is which
	// it has no internal usage
	id ClientID

	logger *log.Logger

	// coder is used for reading from client and writing to the client
	coder packet.MsgCoder

	// reqCh allows getting access to incomming Requests
	// all succesfully read requests from Client will be added to reqCh
	reqCh chan<- Request

	// writeCh keeps track of Messages that should be written to the Client
	writeCh chan packet.Message

	// conn is used to give Client power of closing connection
	// this allows Client to make sure that no goroutines are left running
	conn io.Closer

	// quitCh is internal way to stop writing goroutine
	quitCh chan struct{}

	// muStore is read write lock that protects store
	muStore sync.RWMutex

	// store is key value store can be used to store information about the client
	store map[string]interface{}

	muRoom sync.RWMutex
	room   Room
}

// NewClient should be used from creating a new Client
func NewClient(id ClientID, logger *log.Logger, conn io.Closer, coder packet.MsgCoder, reqCh chan<- Request) *Client {
	return &Client{
		id:      id,
		logger:  logger,
		conn:    conn,
		coder:   coder,
		reqCh:   reqCh,
		writeCh: make(chan packet.Message, 20),
		quitCh:  make(chan struct{}),
		store:   make(map[string]interface{}),
	}
}

// ID returns id
func (c *Client) ID() ClientID {
	return c.id
}

// Start runs startWriter in goroutine and startReader that will block
// It ends once both startWriter and startReader has been closed
func (c *Client) Start() error {
	c.logger.Printf("Client %d starting reading and writing", c.id)
	go c.startWriter()
	c.startReader()

	c.quitCh <- struct{}{}
	c.logger.Printf("Client %d closed", c.id)
	return nil
}

// Write adds Message to queue to be written to the Client
func (c *Client) Write(msg packet.Message) {
	c.writeCh <- msg
}

// Close closes connection which causes read and writing to stop aswell
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// startReader reads from the Client until it receives an error
// error usually caused if connection is closed
// puts incomming Requests to reqCh
func (c *Client) startReader() error {
	req := Request{Client: c, Message: packet.Message{}}

	for {
		err := c.coder.Decode(&req.Message)
		if err != nil {
			return err
		}
		c.logger.Printf("Client %d received message %s", c.id, req.Message)
		c.reqCh <- req
	}
}

// startWriter waits for Messages from on writeCh and then writes them to the Client
// loop will be closed once startReader goroutine is stopped
func (c *Client) startWriter() {
loop:
	for {
		select {
		case msg := <-c.writeCh:
			err := c.coder.Encode(msg)
			if err != nil {
				c.logger.Printf("Client %d writing error %v", c.id, err)
			}
		case <-c.quitCh:
			break loop
		}
	}
}

// SaveValue saves value to key value store
func (c *Client) SaveValue(key string, val interface{}) {
	c.muStore.Lock()
	defer c.muStore.Unlock()
	c.store[key] = val
}

// LoadValue loads value from key value store
func (c *Client) LoadValue(key string) (interface{}, error) {
	c.muStore.RLock()
	defer c.muStore.RUnlock()
	if v, ok := c.store[key]; ok {
		return v, nil
	}
	return nil, errors.New("Key not found " + key)
}

func (c *Client) Room() Room {
	c.muRoom.RLock()
	defer c.muRoom.RUnlock()
	return c.room
}

func (c *Client) setRoom(room Room) {
	c.muRoom.Lock()
	c.room = room
	c.muRoom.Unlock()
}
