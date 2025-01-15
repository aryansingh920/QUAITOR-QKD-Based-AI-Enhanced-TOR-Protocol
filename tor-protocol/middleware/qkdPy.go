package middleware

import (
	"bytes"
	"fmt"

	"os/exec"

	"strings"
)

// In middleware.go (or wherever runPythonCommand is):
func runPythonCommand(mode, message string) (string, error) {
    cmd := exec.Command("python", "../qkd/main.py",
        "--mode", mode,
        "--message", message,
        "--key-length", "256",
    )

    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &out

    err := cmd.Run()
    if err != nil {
        return "", fmt.Errorf("python script error: %w - details: %s", err, out.String())
    }

    // Raw output from Python (which may include newline logs)
    rawOutput := out.String()

    // Replace newlines with spaces (or remove them entirely) to avoid breaking the URL
    // Then trim leading/trailing whitespace
    sanitizedOutput := strings.ReplaceAll(rawOutput, "\n", " ")
    sanitizedOutput = strings.ReplaceAll(sanitizedOutput, "\r", " ")
    sanitizedOutput = strings.TrimSpace(sanitizedOutput)

    return sanitizedOutput, nil
}


func encryptMessage(plaintext string) (string, error) {
    return runPythonCommand("encrypt", plaintext)
}

func decryptMessage(ciphertext string) (string, error) {
    return runPythonCommand("decrypt", ciphertext)
}
