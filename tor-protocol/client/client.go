package client

import (
	"fmt"
	"tor-protocol/protocol"
)

func SendRequest() error {
	// Create a custom header
	header := &protocol.CustomHeader{
		Version:     1,
		RouteID:     12345,
		Timestamp:   1672531200,
		Encrypted:   true,
		PayloadSize: 1024,
	}

	// Serialize the header
	data, err := header.Serialize()
	if err != nil {
		return err
	}

	fmt.Printf("Sending request with header: %x\n", data)
	return nil

	// // Create an HTTP request
	// req, err := http.NewRequest("POST", "http://localhost:8080/secure-route", bytes.NewBuffer([]byte("Test Payload")))
	// if err != nil {
	// 	return err
	// }
	// req.Header.Set("X-Custom-Protocol", string(data))

	// // Send the request
	// resp, err := http.DefaultClient.Do(req)
	// if err != nil {
	// 	return err
	// }
	// defer resp.Body.Close()

	// fmt.Printf("Response: %s\n", resp.Status)
	// return nil
}
