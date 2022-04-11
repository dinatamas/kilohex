package terminal

import (
    "fmt"
    "io"
    "log"
)

// Provides the most common terminal-drawing functionalities.
type Window struct {
    H, W int
    Buffer string
    OriginalConfig TerminalConfig
    OriginalLogWriter io.Writer
}

// Emits screen clearing Ansi escape codes.
// Repositions the cursor to the home (top left) position.
// Note:
//   - AnsiClearScreen places the current screen in the terminal scrollback buffer.
//   - AnsiClearBelow clears the screen just the same, without retaining history.
func (w *Window) Clear() {
    w.Buffer += AnsiCursorHome + AnsiClearBelow
}

// Synchronize the window's buffer to the terminal screen.
func (w *Window) Flush() {
    fmt.Print(w.Buffer)
    w.Buffer = ""
}

// Creates a new window and configures the terminal to work with it.
// The Restore() method of the returned window should be deferred!
func NewWindow() Window {
    config := GetConfig()
    SetRawMode()
    rows, cols := GetSize()
    writer := log.Writer()
    log.SetOutput(RawLogWriter{})
    w := Window{rows, cols, "", config, writer}
    // TODO: The original terminal content should be restored!
    // Use save screen or alternative buffer private sequences?
    w.Buffer += AnsiClearScreen + AnsiCursorHome
    return w
}

// Prints the string starting from the current cursor position.
// The cursor will be displaced accordingly.
func (w *Window) Print(s string) {
    w.Buffer += s
}

// Prints the string starting from the specified location.
// The cursor will remain in its original location.
func (w *Window) PrintAt(s string, y int, x int) {
    w.Buffer += AnsiCursorSave
    w.Buffer += fmt.Sprintf(AnsiCursorSetPosition, y + 1, x + 1)
    w.Buffer += s
    w.Buffer += AnsiCursorRestore

}

// Retrieves terminal dimensions and updates the window's size accordingly.
// Returns true if there was a change, false otherwise.
func (w *Window) DetectResize() bool {
    rows, cols := GetSize()
    if rows != w.H || cols != w.W {
        w.H, w.W = rows, cols
        return true
    }
    return false
}

// Restores the terminal to its original configuration.
func (w *Window) Restore() {
    w.Clear()
    w.Flush()
    SetConfig(w.OriginalConfig)
    log.SetOutput(w.OriginalLogWriter)
}

// Sets the window's cursor position.
func (w *Window) SetCursorPosition(y int, x int) {
    w.Buffer += fmt.Sprintf(AnsiCursorSetPosition, y + 1, x + 1)
}
