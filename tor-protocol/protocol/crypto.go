package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// These are placeholder "serialization" and "deserialization" methods
// for the RelayCell. In a real Tor network, you'd have complex layering
// and encryption.

func (rc *RelayCell) Serialize() ([]byte, error) {
    // For demonstration:
    // 4 bytes for nextAddr length, nextAddr (string)
    // 4 bytes for payload length, payload (bytes)
    nextAddrBytes := []byte(rc.NextAddr)
    nextLen := int32(len(nextAddrBytes))
    payloadLen := int32(len(rc.Payload))

    buf := new(bytes.Buffer)
    if err := binary.Write(buf, binary.BigEndian, nextLen); err != nil {
        return nil, err
    }
    if nextLen > 0 {
        if _, err := buf.Write(nextAddrBytes); err != nil {
            return nil, err
        }
    }

    if err := binary.Write(buf, binary.BigEndian, payloadLen); err != nil {
        return nil, err
    }
    if payloadLen > 0 {
        if _, err := buf.Write(rc.Payload); err != nil {
            return nil, err
        }
    }
    return buf.Bytes(), nil
}

func DeserializeCell(data []byte) (*RelayCell, error) {
    buf := bytes.NewReader(data)

    var nextLen int32
    if err := binary.Read(buf, binary.BigEndian, &nextLen); err != nil {
        return nil, fmt.Errorf("failed to read nextLen: %v", err)
    }

    nextAddrBytes := make([]byte, nextLen)
    if nextLen > 0 {
        if _, err := buf.Read(nextAddrBytes); err != nil {
            return nil, fmt.Errorf("failed to read nextAddrBytes: %v", err)
        }
    }

    var payloadLen int32
    if err := binary.Read(buf, binary.BigEndian, &payloadLen); err != nil {
        return nil, fmt.Errorf("failed to read payloadLen: %v", err)
    }

    payloadBytes := make([]byte, payloadLen)
    if payloadLen > 0 {
        if _, err := buf.Read(payloadBytes); err != nil {
            return nil, fmt.Errorf("failed to read payloadBytes: %v", err)
        }
    }

    return &RelayCell{
        NextAddr: string(nextAddrBytes),
        Payload:  payloadBytes,
    }, nil
}
