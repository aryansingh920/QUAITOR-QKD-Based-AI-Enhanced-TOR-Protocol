/*
Created on 11/01/2025

@author: Aryan

Filename: crypto.go

Relative Path: tor-protocol/protocol/crypto.go
*/

package protocol

import (
	"bytes"
	"crypto/rand"
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

// GenerateEphemeralKeyPair creates a dummy ephemeral key pair for demonstration.
// In a real implementation, you'd use something like X25519, curve25519, or RSA, etc.
func GenerateEphemeralKeyPair() ([]byte, []byte, error) {
    // For demonstration, just generate 32 random bytes for each "key".
    public := make([]byte, 32)
    private := make([]byte, 32)

    _, errPub := rand.Read(public)
    _, errPriv := rand.Read(private)

    if errPub != nil || errPriv != nil {
        return nil, nil, fmt.Errorf("error generating ephemeral key material")
    }

    return public, private, nil
}

// EncryptPayload would symmetrically or asymmetrically encrypt data in a real Tor system.
// We just provide a placeholder here.
func EncryptPayload(payload, key []byte) []byte {
    // Placeholder encryption: in practice you'd do real encryption using ephemeralPublicKey
    // For demonstration, just return payload as-is
    return payload
}

// DecryptPayload would symmetrically or asymmetrically decrypt data in a real Tor system.
// We just provide a placeholder here.
func DecryptPayload(ciphertext, key []byte) []byte {
    // Placeholder decryption: in practice you'd do real decryption using ephemeralPrivateKey
    return ciphertext
}
