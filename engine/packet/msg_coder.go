package packet

import (
	"bufio"
	"io"
)

// MsgCoder provides interface for Encoding and Decoding Message's
type MsgCoder interface {
	Encode(msg Message) error
	Decode(msg *Message) error
}

// MsgCode is basic implementation of MsgCoder interface
type MsgCode struct {
	reader *bufio.Reader
	writer *bufio.Writer
}

// NewMsgCode helper function to create *MsgCode
func NewMsgCode(r io.Reader, w io.Writer) *MsgCode {
	return &MsgCode{
		reader: bufio.NewReader(r),
		writer: bufio.NewWriter(w),
	}
}

// Encode encodes Message to writer
// Example of raw bytes
// MOVE\n
// UID:abc;TOKEN:1231\n
// {"x": 1, "y": 2}\n
// Line break is required for every line, other than that it can be empty
func (m *MsgCode) Encode(msg Message) error {
	cmd := []byte(msg.CMD)
	opts := msg.RawOptions
	payload := msg.Payload

	// Add line breaks
	if len(cmd) == 0 || cmd[len(cmd)-1] != '\n' {
		cmd = append(cmd, '\n')
	}
	if len(opts) == 0 || opts[len(opts)-1] != '\n' {
		opts = append(opts, '\n')
	}
	if len(payload) == 0 || payload[len(payload)-1] != '\n' {
		payload = append(payload, '\n')
	}

	_, err := m.writer.Write(cmd)
	if err != nil {
		return err
	}
	_, err = m.writer.Write(opts)
	if err != nil {
		return err
	}
	_, err = m.writer.Write(payload)
	if err != nil {
		return err
	}
	return m.writer.Flush()
}

// Decode decodes *Message from reader
// Incomming data format should be same as on Encode
func (m *MsgCode) Decode(msg *Message) error {
	cmd, err := m.reader.ReadBytes('\n')
	if err != nil {
		return err
	}
	opt, err := m.reader.ReadBytes('\n')
	if err != nil {
		return err
	}
	payload, err := m.reader.ReadBytes('\n')
	if err != nil {
		return err
	}

	// Remove trailing line breaks
	if len(cmd) != 0 && cmd[len(cmd)-1] == '\n' {
		cmd = cmd[:len(cmd)-1]
	}
	if len(opt) != 0 && opt[len(opt)-1] == '\n' {
		opt = opt[:len(opt)-1]
	}
	if len(payload) != 0 && payload[len(payload)-1] == '\n' {
		payload = payload[:len(payload)-1]
	}

	msg.CMD = Command(cmd)
	msg.RawOptions = opt
	msg.Payload = payload

	return nil
}
