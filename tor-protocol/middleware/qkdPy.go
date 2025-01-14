package middleware

import (
	"bytes"
	"fmt"

	"os/exec"

	"strings"
)



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

    // Trim to avoid trailing newlines
    return strings.TrimSpace(out.String()), nil
}

func encryptMessage(plaintext string) (string, error) {
    return runPythonCommand("encrypt", plaintext)
}

func decryptMessage(ciphertext string) (string, error) {
    return runPythonCommand("decrypt", ciphertext)
}
