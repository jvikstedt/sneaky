package engine

import "github.com/jvikstedt/sneaky/engine/packet"

// Request is used for keeping track of Client and Message
type Request struct {
	Client *Client
	packet.Message
}
