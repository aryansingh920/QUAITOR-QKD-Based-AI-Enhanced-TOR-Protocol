/*
Updated on 11/01/2025
@author: Aryan

Placeholder cryptographic routines + RelayCell serialization.
*/
package protocol

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
)

// -- RelayCell Serialization / Deserialization --

func (rc *RelayCell) Serialize() ([]byte, error) {
    // We'll store:
    //   - PrevAddrLen + PrevAddr
    //   - NextAddrLen + NextAddr
    //   - PayloadLen  + Payload
    //   - 2 bools for IsExitRequest, IsExitResponse

    prevAddrBytes := []byte(rc.PrevAddr)
    nextAddrBytes := []byte(rc.NextAddr)

    buf := new(bytes.Buffer)
    // write prevAddr
    if err := binary.Write(buf, binary.BigEndian, int32(len(prevAddrBytes))); err != nil {
        return nil, err
    }
    buf.Write(prevAddrBytes)
    // write nextAddr
    if err := binary.Write(buf, binary.BigEndian, int32(len(nextAddrBytes))); err != nil {
        return nil, err
    }
    buf.Write(nextAddrBytes)
    // write payload
    if err := binary.Write(buf, binary.BigEndian, int32(len(rc.Payload))); err != nil {
        return nil, err
    }
    buf.Write(rc.Payload)
    // write booleans
    var exitReq, exitRes int8
    if rc.IsExitRequest {
        exitReq = 1
    }
    if rc.IsExitResponse {
        exitRes = 1
    }
    if err := binary.Write(buf, binary.BigEndian, exitReq); err != nil {
        return nil, err
    }
    if err := binary.Write(buf, binary.BigEndian, exitRes); err != nil {
        return nil, err
    }
    return buf.Bytes(), nil
}

func ParseRelayCell(data []byte) (*RelayCell, error) {
    buf := bytes.NewReader(data)

    var prevAddrLen int32
    if err := binary.Read(buf, binary.BigEndian, &prevAddrLen); err != nil {
        return nil, fmt.Errorf("failed to read prevAddrLen: %v", err)
    }
    prevAddrBytes := make([]byte, prevAddrLen)
    if _, err := buf.Read(prevAddrBytes); err != nil {
        return nil, fmt.Errorf("failed to read prevAddr: %v", err)
    }

    var nextAddrLen int32
    if err := binary.Read(buf, binary.BigEndian, &nextAddrLen); err != nil {
        return nil, fmt.Errorf("failed to read nextAddrLen: %v", err)
    }
    nextAddrBytes := make([]byte, nextAddrLen)
    if _, err := buf.Read(nextAddrBytes); err != nil {
        return nil, fmt.Errorf("failed to read nextAddr: %v", err)
    }

    var payloadLen int32
    if err := binary.Read(buf, binary.BigEndian, &payloadLen); err != nil {
        return nil, fmt.Errorf("failed to read payloadLen: %v", err)
    }
    payloadBytes := make([]byte, payloadLen)
    if _, err := buf.Read(payloadBytes); err != nil {
        return nil, fmt.Errorf("failed to read payload: %v", err)
    }

    var exitReq, exitRes int8
    if err := binary.Read(buf, binary.BigEndian, &exitReq); err != nil {
        return nil, fmt.Errorf("failed to read IsExitRequest: %v", err)
    }
    if err := binary.Read(buf, binary.BigEndian, &exitRes); err != nil {
        return nil, fmt.Errorf("failed to read IsExitResponse: %v", err)
    }

    rc := &RelayCell{
        PrevAddr:       string(prevAddrBytes),
        NextAddr:       string(nextAddrBytes),
        Payload:        payloadBytes,
        IsExitRequest:  (exitReq == 1),
        IsExitResponse: (exitRes == 1),
    }
    return rc, nil
}

// -- Dummy ephemeral key generation --

func GenerateEphemeralKeyPair() ([]byte, []byte, error) {
    pub := make([]byte, 32)
    priv := make([]byte, 32)
    if _, err := rand.Read(pub); err != nil {
        return nil, nil, fmt.Errorf("error generating public key: %v", err)
    }
    if _, err := rand.Read(priv); err != nil {
        return nil, nil, fmt.Errorf("error generating private key: %v", err)
    }
    return pub, priv, nil
}
