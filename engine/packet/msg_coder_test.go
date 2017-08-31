package packet_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/jvikstedt/sneaky/engine/packet"
)

var messages = []packet.Message{
	{CMD: "JUMP", RawOptions: []byte(""), Payload: []byte("")},
	{CMD: "MOVE", RawOptions: []byte(""), Payload: []byte(`{"x": 1, "y": 2}`)},
	{CMD: "DISCONNECT", RawOptions: []byte("UID:abc123"), Payload: []byte("")},
}

func TestEncode(t *testing.T) {
	b := &bytes.Buffer{}

	msgCode := packet.NewMsgCode(b, b)

	for _, m := range messages {
		err := msgCode.Encode(m)
		if err != nil {
			t.Error(err)
		}

		expected := fmt.Sprintf("%s\n%s\n%s\n", m.CMD, string(m.RawOptions), string(m.Payload))
		actual := string(b.Bytes())
		if actual != expected {
			t.Errorf("Expected %s but got %s", expected, actual)
		}
		b.Reset()
	}
}

func TestDecode(t *testing.T) {
	b := &bytes.Buffer{}

	msgCode := packet.NewMsgCode(b, b)

	for _, m := range messages {
		data := fmt.Sprintf("%s\n%s\n%s\n", m.CMD, string(m.RawOptions), string(m.Payload))
		actual := packet.Message{}

		_, err := b.Write([]byte(data))
		if err != nil {
			t.Error(err)
		}
		err = msgCode.Decode(&actual)
		if err != nil {
			t.Error(err)
		}

		if actual.CMD != m.CMD {
			t.Errorf("Expected %s but got %s", m.CMD, actual.CMD)
		}

		if string(actual.RawOptions) != string(m.RawOptions) {
			t.Errorf("Expected %s but got %s", string(m.CMD), string(actual.CMD))
		}

		if string(actual.Payload) != string(m.Payload) {
			t.Errorf("Expected %s but got %s", string(m.CMD), string(actual.CMD))
		}
	}
}
