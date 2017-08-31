package packet

import (
	"encoding/json"
	"strings"
)

type Command string

// Message is struct format of incomming Data
type Message struct {
	// CMD (short for command) is usually just short string like MOVE or JUMP
	CMD Command

	// RawOptions is raw byte version of options
	RawOptions []byte

	// Options is RawOptions parsed to key value format
	Options map[string]string

	// Payload is anything that is required for CMD excecuting
	Payload []byte
}

// ParseOptions Serializes RawOptions to Options
func (msg *Message) ParseOptions() {
	msg.Options = make(map[string]string)

	for _, v := range strings.Split(string(msg.RawOptions), ";") {
		keyAndValue := strings.Split(v, ":")
		msg.Options[keyAndValue[0]] = keyAndValue[1]
	}
}

func NewJsonMessage(cmd Command, v interface{}, options map[string]string) Message {
	payload, _ := json.Marshal(v)

	return Message{
		CMD:     cmd,
		Payload: payload,
		Options: options,
	}
}
