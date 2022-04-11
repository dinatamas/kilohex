package terminal

import (
    "fmt"
)

// Use this writer to log to the terminal in raw mode.
type RawLogWriter struct {}

func (_ RawLogWriter) Write(p []byte) (n int, err error) {
    return fmt.Print(append(p, '\r'))
}
