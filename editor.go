package main

import (
    "strings"

    "github.com/dinatamas/kilohex/terminal"
)

const (
    STATE_CLEAN = 0
    STATE_DIRTY = 1
    STATE_SAVING = 2
    STATE_EXITED = 3
)

type Editor struct {
    State    int
    Filename string
    Buffer   string
    Rows     []string
    W        terminal.Window
    Cy, Cx   int
    Oy, Ox   int
    LineEnding string
}

func (e *Editor) Restore() {
    e.W.Restore()
}

func (e *Editor) Draw() {
    e.W.DetectResize()
    e.W.Clear()

    // Draw rows.
    for i := 0; i < e.W.H - 1; i++ {
        if e.Oy + i >= len(e.Rows) { break; }
        if e.Ox > len(e.Rows[e.Oy + i]) { continue; }
        lastX := e.Ox + e.W.W
        if lastX > len(e.Rows[e.Oy + i]) {
            lastX = len(e.Rows[e.Oy + i])
        }
        e.W.Print(e.Rows[e.Oy + i][e.Ox:lastX] + "\r\n")
    }

    // Draw status bar.
    status := e.Filename
    switch e.State {
    case STATE_DIRTY:
        status += "*"
    case STATE_SAVING:
        status += " - Save?"
    }
    status += strings.Repeat(" ", e.W.W - len(status))
    e.W.PrintAt(terminal.AnsiInverse + status + terminal.AnsiInverseReset, e.W.H, 0)

    e.W.SetCursorPosition(e.Cy - e.Oy, e.Cx)
    e.W.Flush()
}

// Helper function for HandleKey to avoid code duplication.
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
    // Move offsets if required.
    if e.Cy < e.Oy {
        e.Oy = e.Cy
    } else if e.Cy > e.Oy + e.W.H - 2 {
        e.Oy = e.Cy - e.W.H + 2
    }
    if e.Cx < e.Ox {
        e.Ox = e.Cx
    } else if e.Cx > e.Ox + e.W.W - 1 {
        e.Ox = e.Cx - e.W.W + 1
    }
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

func NewEditor(filename string, buffer string) Editor {
    lineEnding := "\n"
lineEndingLoop:
    for _, c := range buffer {
        switch c {
        case '\r':
            lineEnding = "\r"
            break lineEndingLoop
        case '\n':
            break lineEndingLoop
        }
    }
    lines := strings.Split(buffer, lineEnding)
    window := terminal.NewWindow()
    editor := Editor{
        State: STATE_CLEAN,
        Filename: filename,
        Buffer: buffer,
        Rows: lines,
        W: window,
        Cy: 0, Cx: 0,
        Oy: 0, Ox: 0,
        LineEnding: lineEnding,
    }
    return editor
}

func (e *Editor) Run() string {
    for e.State != STATE_EXITED {
        e.Draw()
        e.HandleKey(terminal.ReadKey())
    }
    return e.Buffer
}

func (e *Editor) Save() {
    e.State = STATE_CLEAN
    e.Buffer = strings.Join(e.Rows, e.LineEnding)
}
