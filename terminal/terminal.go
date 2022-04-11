package terminal

import (
    "fmt"
)

func GetCursorPosition() (y int, x int) {
    fmt.Print(AnsiCursorGetPosition)
    // The terminal will write the cursor position to stdin.
    fmt.Scanf(AnsiCursorPositionReport, &y, &x)
    return y, x
}

func SetCursorPosition(y int, x int) {
    fmt.Printf(AnsiCursorSetPosition, y, x)
}
