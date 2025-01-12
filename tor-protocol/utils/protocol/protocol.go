package protocol

import (
	"bytes"
	"encoding/binary"
)

// CustomHeader represents the protocol header structure
type CustomHeader struct {
	Version     uint8  // Protocol version
	RouteID     uint16 // Unique ID for the route
	Timestamp   int64  // Unix timestamp
	Encrypted   bool   // Encryption flag
	PayloadSize uint32 // Size of the payload
}

// Serialize serializes the header into a byte slice
func (h *CustomHeader) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, h); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Deserialize deserializes a byte slice into a CustomHeader
func Deserialize(data []byte) (*CustomHeader, error) {
	h := &CustomHeader{}
	buf := bytes.NewReader(data)
	if err := binary.Read(buf, binary.BigEndian, h); err != nil {
		return nil, err
	}
	return h, nil
}
