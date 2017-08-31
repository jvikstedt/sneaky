package packet

import (
	"testing"
)

func TestRequest(t *testing.T) {
	msg := Message{RawOptions: []byte(`FIRST:123;SECOND:Something`)}
	msg.ParseOptions()

	if msg.Options["FIRST"] != "123" {
		t.Errorf("Expected %s but got %s", "123", msg.Options["FIRST"])
	}

	if msg.Options["SECOND"] != "Something" {
		t.Errorf("Expected %s but got %s", "Something", msg.Options["SECOND"])
	}
}
