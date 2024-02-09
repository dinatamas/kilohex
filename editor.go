package main

import (
    "fmt"
    "strings"

    "github.com/dinatamas/kilohex/terminal"
)

const (
    STATE_CLEAN  = 0
    STATE_DIRTY  = 1
    STATE_SAVING = 2
    STATE_EXITED = 3
)

//===================================================================
//
// External Editor methods.
//
//===================================================================

type Editor struct {
    State      int
    Filename   string
    Buffer     []byte
    W          terminal.Window
    Cy, Cx     int
    Oy, Ox     int
    Error      error
}

func NewEditor(filename string, buffer []byte) Editor {
    // Instantiate the window and the editor.
    window := terminal.NewWindow()
    editor := Editor{
        State: STATE_CLEAN,
        Filename: filename,
        Buffer: buffer,
        W: window,
        Cy: 0, Cx: 0,
        Oy: 0, Ox: 0,
    }
    return editor
}

func (e *Editor) Close() {
    e.W.Restore()
}

func (e *Editor) Run() ([]byte, error) {
    // Infinite event loop.
    for {
        if e.State == STATE_EXITED { return e.Buffer, nil }

        e.Draw()
        e.HandleKey(terminal.ReadKey())
    }
}

func (e *Editor) Save() {}

//===================================================================
//
// Physical and logical calculated cursor positions.
//
//===================================================================

// Logical cursor Y-position (relative to buffer).
func (e *Editor) lCy() int { return e.Cy }
// Logical cursor X-position (relative to logical line).
func (e *Editor) lCx() int { return 10 + e.Cx * 3 + (e.Cx / 8) }

// Physical cursor Y-position (relative to terminal).
func (e *Editor) tCy() int { return e.lCy() - e.Oy }
// Phyiscal cursor X-position (relative to terminal).
func (e *Editor) tCx() int { return e.lCx() - e.Ox }

//===================================================================
//
// Physical and logical boundaries for the displayed content.
//
//===================================================================

// idx of first / last+1 physical row of content area (relative to terminal)
func (e *Editor) tMinY() int { return 0 }
func (e *Editor) tMaxY() int { return e.W.H - 1 }
// idx of first / last+1 line in the entire buffer (relative to buffer)
func (e *Editor) bMinY() int { return 0 }
func (e *Editor) bMaxY() int { return (len(e.Buffer) + 16 - 1) / 16 }
// idx of first / last+1 logical line in content area (relative to buffer)
func (e *Editor) lMinY() int { return e.Oy - e.tMinY() }
func (e *Editor) lMaxY() int { return min(e.lMinY() + e.tMaxY(), e.bMaxY()) }

// idx of first / last+1 physical col of content area (relative to terminal)
func (e *Editor) tMinX(ly int) int { return 0 }
func (e *Editor) tMaxX(ly int) int { return e.W.W }
// idx of first / last+1 byte in the logical line (relative to logical line)
func (e *Editor) bMinX(ly int) int { return 0 }
func (e *Editor) bMaxX(ly int) int { return min(len(e.Buffer) - ly * 16, 16) }
// idx of first / last+1 byte in content area (relative to logical line)
func (e *Editor) lMinX(ly int) int { return e.Ox - e.tMinX(ly) }
func (e *Editor) lMaxX(ly int) int { return min(e.lMinX(ly) + e.tMaxX(ly), 78) }

//===================================================================
//
// Internal Editor methods.
//
//===================================================================

func (e *Editor) Draw() {
    e.W.DetectResize()
    e.W.Clear()

    // Iterate each logical line with content.
    for ly := e.lMinY(); ly < e.lMaxY(); ly++ {
        bytes := e.Buffer[ly * 16 + e.bMinX(ly):ly * 16 + e.bMaxX(ly)]

        // Construct the logical content line.
        var loLine, offset, encoded, decoded string
        offset = fmt.Sprintf("%08x", ly * 16)
        for i, b := range bytes {
            encoded += fmt.Sprintf("%02x ", b)
            if i == 7 { encoded += " " }
        }
        encoded = fmt.Sprintf("%-49s", encoded)
        decoded = "................"
        loLine = offset + "  " + encoded + " |" + decoded + "|"

        // Construct the physical displayed line.
        phLine := loLine[e.Ox:e.lMaxX(ly)]
        e.W.Print(phLine + "\r\n")
    }

    // Draw status bar.
    status := e.Filename
    switch e.State {

        case STATE_DIRTY:
            status += "*"

        case STATE_SAVING:
            status += " - Save?"
    }
    if e.W.W > len(status) {
      status += strings.Repeat(" ", e.W.W - len(status))
    }
    e.W.PrintAt(terminal.AnsiInverse + status[0:e.W.W] + terminal.AnsiInverseReset, e.W.H, 0)

    e.W.SetCursorPosition(e.tCy(), e.tCx())
    e.W.Flush()
}

func (e *Editor) HandleKey(key terminal.Key) {
    switch e.State {

        case STATE_CLEAN:
            switch key.Value {
                case terminal.KeyCtrlC:
                    e.State = STATE_EXITED

                default:
                    e.handleEditorKey(key)
            }

        case STATE_DIRTY:
            switch key.Value {
                case terminal.KeyCtrlC:
                    e.State = STATE_SAVING

                case terminal.KeyCtrlS:
                    e.Save()

                default:
                    e.handleEditorKey(key)
            }

        case STATE_SAVING:
            switch key.Value {
                case terminal.KeyCtrlC:
                    e.State = STATE_EXITED

                case terminal.KeyEsc:
                    e.State = STATE_DIRTY

                case terminal.KeyEnter:
                    e.Save()
                    e.State = STATE_EXITED
            }
    }
}

func (e *Editor) handleEditorKey(key terminal.Key) {
    if key.Type == terminal.KEY_RUNE {

        e.State = STATE_DIRTY

    } else {
        switch key.Value {

            case terminal.KeyUp:
                if e.Cy > e.bMinY() {
                    e.Cy--
                    if e.Cx > e.bMaxX(e.Cy) - 1 {
                        e.Cx = e.bMaxX(e.Cy) - 1
                    }
                }

            case terminal.KeyDown:
                if e.Cy < e.bMaxY() - 1 {
                    e.Cy++
                    if e.Cx > e.bMaxX(e.Cy) - 1 {
                        e.Cx = e.bMaxX(e.Cy) - 1
                    }
                }

            case terminal.KeyLeft:
                if e.Cx > e.bMinX(e.Cy) {
                    e.Cx--
                }

            case terminal.KeyRight:
                if e.Cx < e.bMaxX(e.Cy) - 1 {
                    e.Cx++
                }

            case terminal.KeyHome:
                e.Cx = e.bMinX(e.Cy)

            case terminal.KeyEnd:
                e.Cx = e.bMaxX(e.Cy) - 1
        }
    }

    // Adjust vertical offsets.
    if e.lCy() < e.lMinY() {
        e.Oy = e.lCy() - e.tMinY()
    } else if e.lCy() > e.lMaxY() - 1 {
        e.Oy = e.lCy() - (e.tMaxY() - 1)
    }

    // Adjust horizontal offsets.
    if e.lCx() < e.lMinX(e.Cy) {
        e.Ox = e.lCx() - e.tMinX(e.Cy)
    } else if e.lCx() > e.lMaxX(e.Cy) - 1 {
        // Add 1 character scrolloff to not cut bytes in half.
        e.Ox = e.lCx() - (e.tMaxX(e.Cy) - 1) + 1
    }
}
