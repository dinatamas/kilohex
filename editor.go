package main

import (
    "encoding/hex"
    "fmt"
    "strings"

    "github.com/dinatamas/kilohex/terminal"
)

const (
    STATE_CLEAN  = 0
    STATE_DIRTY  = 1
    STATE_SAVING = 2
    STATE_EXITED = 3
    STATE_ERROR  = 4
)

type Editor struct {
    State      int
    Filename   string
    Buffer     []byte           // Saved continuous file content.
    Rows       []string         // Unsaved editor lines.
    W          terminal.Window
    Cy, Cx     int              // Cursor offset compared to screen.
    Oy, Ox     int              // Screen offset compared to file.
    Error      error
}

func NewEditor(filename string, buffer []byte) Editor {
    // Represent the buffer as rows of hexadecimal strings.
    hexbuf := hex.EncodeToString(buffer)
    var lines []string
    for i := 0; i < len(hexbuf); i += 32 {
        j := i + 32
        if j > len(hexbuf) { j = len(hexbuf) }
        lines = append(lines, hexbuf[i:j])
    }

    // Instantiate the window and the editor.
    window := terminal.NewWindow()
    editor := Editor{
        State: STATE_CLEAN,
        Filename: filename,
        Buffer: buffer,
        Rows: lines,
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
        if e.State == STATE_ERROR  { return []byte{}, e.Error }

        e.Draw()
        e.HandleKey(terminal.ReadKey())
    }
}

func (e *Editor) Save() {
    // Reconstruct the buffer from the hex strings.
    hexbuf := strings.Join(e.Rows, "")
    buffer, err := hex.DecodeString(hexbuf)

    if err == nil {
        // Save successful.
        e.Buffer = buffer
        e.State = STATE_CLEAN
    } else {
        // Save failed.
        e.Error = err
        e.State = STATE_ERROR
    }
}

func (e *Editor) Draw() {
    e.W.DetectResize()
    e.W.Clear()

    // Draw rows.
    for i := 0; i < e.W.H - 1; i++ {

        if e.Oy + i >= len(e.Rows) {
            break;
        }
        if e.Ox > len(e.Rows[e.Oy + i]) {
            continue;
        }

        // Create the displayed row with left and right sides.
        hexpos := fmt.Sprintf("%08x", (e.Oy + i) * 16)
        content := ""
        for j := 0; j < 32; j += 1 {
            if j < len(e.Rows[e.Oy + i]) {
                content += string(e.Rows[e.Oy + i][j])
            } else {
                content += " "
            }
            if j % 2 == 1 {
                content += " "
            }
            if j == 15 {
                content += " "
            }
        }
        decoded := ""

        loLine := hexpos + "  " + content + " |" + decoded + "|\r\n"
        e.W.Print(loLine)
        continue

        lastX := e.Ox + e.W.W
        if lastX > len(e.Rows[e.Oy + i]) {
            lastX = len(e.Rows[e.Oy + i])
        }

        octets := e.Rows[e.Oy + i][e.Ox:lastX]
        for j := e.Ox; j < lastX; j += 1 {
            if j > 0 && j % 2 == 0 {
                e.W.Print(" ")
            }
        }

        var bytes []string
        for j := e.Ox; j < lastX; j += 2 {
            k := j + 2
            if k > lastX { k = lastX }
            bytes = append(bytes, octets[j:k])
        }
        e.W.Print(strings.Join(bytes, " ") + "\r\n")
    }

    // Draw status bar.
    status := e.Filename
    switch e.State {
        case STATE_DIRTY:
            status += "*"

        case STATE_ERROR:
            status += "!"

        case STATE_SAVING:
            status += " - Save?"
    }
    status += strings.Repeat(" ", e.W.W - len(status))
    e.W.PrintAt(terminal.AnsiInverse + status + terminal.AnsiInverseReset, e.W.H, 0)

    e.W.SetCursorPosition(e.Cy - e.Oy, e.Cx)
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
        e.Rows[e.Cy] = e.Rows[e.Cy][:e.Cx] + key.Value + e.Rows[e.Cy][e.Cx:]
        e.Cx++
    } else {
        switch key.Value {
            case terminal.KeyUp:
                if e.Cy > 0 {
                    e.Cy--
                    if e.Cx > len(e.Rows[e.Cy]) {
                        e.Cx = len(e.Rows[e.Cy])
                    }
                }

            case terminal.KeyDown:
                if e.Cy < len(e.Rows) - 1 {
                    e.Cy++
                    if e.Cx > len(e.Rows[e.Cy]) {
                        e.Cx = len(e.Rows[e.Cy])
                    }
                }

            case terminal.KeyRight:
                if e.Cx < len(e.Rows[e.Cy]) {
                    e.Cx++
                }

            case terminal.KeyLeft:
                if e.Cx > 0 {
                    e.Cx--
                }

            case terminal.KeyBackspace, terminal.Ascii_BS:
                if e.Cx > 0 {
                    e.State = STATE_DIRTY
                    e.Rows[e.Cy] = e.Rows[e.Cy][:e.Cx-1] + e.Rows[e.Cy][e.Cx:]
                    e.Cx--
                } else if len(e.Rows[e.Cy]) == 0 && e.Cy > 0 {
                    e.State = STATE_DIRTY
                    e.Rows = append(e.Rows[:e.Cy], e.Rows[e.Cy+1:]...)
                    e.Cy--
                    e.Cx = len(e.Rows[e.Cy])
                }

            case terminal.KeyEnter:
                e.State = STATE_DIRTY
                e.Rows = append(e.Rows[:e.Cy+1], e.Rows[e.Cy:]...)
                e.Rows[e.Cy+1] = e.Rows[e.Cy][e.Cx:]
                e.Rows[e.Cy] = e.Rows[e.Cy][:e.Cx]
                e.Cy++
                e.Cx = 0

            case terminal.KeyDelete:
                if e.Cx == len(e.Rows[e.Cy]) && e.Cy < len(e.Rows) - 1 {
                    e.State = STATE_DIRTY
                    e.Rows[e.Cy] = e.Rows[e.Cy] + e.Rows[e.Cy+1]
                    e.Rows = append(e.Rows[:e.Cy+1], e.Rows[e.Cy+2:]...)
                } else if e.Cx < len(e.Rows[e.Cy]) {
                    e.State = STATE_DIRTY
                    e.Rows[e.Cy] = e.Rows[e.Cy][:e.Cx] + e.Rows[e.Cy][e.Cx+1:]
                }

            case terminal.KeyHome:
                e.Cx = 0

            case terminal.KeyEnd:
                e.Cx = len(e.Rows[e.Cy])

            case terminal.KeyTab:
                break // TODO

            case terminal.KeyPgUp:
                break // TODO

            case terminal.KeyPgDn:
                break // TODO
        }
    }

    // Adjust vertical offsets.
    if e.Cy < e.Oy {
        e.Oy = e.Cy
    } else if e.Cy > e.Oy + e.W.H - 2 {
        e.Oy = e.Cy - e.W.H + 2
    }

    // Adjust horizontal offsets.
    if e.Cx < e.Ox {
        e.Ox = e.Cx
    } else if e.Cx > e.Ox + e.W.W - 1 {
        e.Ox = e.Cx - e.W.W + 1
    }
}
