package terminal

import (
    "bufio"
    "os"
    "unicode"
)

// Key types returned by ReadKey().
const (
    KEY_ANSI = 0
    KEY_CTRL = 1
    KEY_RUNE = 2
    KEY_ALT  = 3
)

// ASCII control characters.
const (
    // Named physical keys.
    KeyTab       = "\x09"
    KeyEnter     = "\x0d"
    KeyEsc       = "\x1b"
    KeySpace     = "\x20"
    KeyBackspace = "\x7f"  // Ascii_BS is also worth checking!

    // Ctrl modifiers.
                       Ascii_NUL = "\x00" // null
    KeyCtrlA = "\x01"; Ascii_SOH = "\x01" // start of heading
    KeyCtrlB = "\x02"; Ascii_STX = "\x02" // start of text
    KeyCtrlC = "\x03"; Ascii_ETX = "\x03" // end of text
    KeyCtrlD = "\x04"; Ascii_EOT = "\x04" // end of transmission
    KeyCtrlE = "\x05"; Ascii_ENQ = "\x05" // enquiry
    KeyCtrlF = "\x06"; Ascii_ACK = "\x06" // acknowledge
    KeyCtrlG = "\x07"; Ascii_BEL = "\x07" // bell
    KeyCtrlH = "\x08"; Ascii_BS  = "\x08" // backspace
    KeyCtrlI = "\x09"; Ascii_TAB = "\x09" // horizontal tab
    KeyCtrlJ = "\x0a"; Ascii_LF  = "\x0a" // line feed (new line)
    KeyCtrlK = "\x0b"; Ascii_VT  = "\x0b" // vertical tab
    KeyCtrlL = "\x0c"; Ascii_FF  = "\x0c" // form feed (new page)
    KeyCtrlM = "\x0d"; Ascii_CR  = "\x0d" // carriage return
    KeyCtrlN = "\x0e"; Ascii_SO  = "\x0e" // shift out
    KeyCtrlO = "\x0f"; Ascii_SI  = "\x0f" // shift in
    KeyCtrlP = "\x10"; Ascii_DLE = "\x10" // data link escape
    KeyCtrlQ = "\x11"; Ascii_DC1 = "\x11" // device control 1
    KeyCtrlR = "\x12"; Ascii_DC2 = "\x12" // device control 2
    KeyCtrlS = "\x13"; Ascii_DC3 = "\x13" // device control 3
    KeyCtrlT = "\x14"; Ascii_DC4 = "\x14" // device control 4
    KeyCtrlU = "\x15"; Ascii_NAK = "\x15" // negative acknowledge
    KeyCtrlV = "\x16"; Ascii_SYN = "\x16" // synchronous idle
    KeyCtrlW = "\x17"; Ascii_ETB = "\x17" // end of transmission block
    KeyCtrlX = "\x18"; Ascii_CAN = "\x18" // cancel
    KeyCtrlY = "\x19"; Ascii_EM  = "\x19" // end of medium
    KeyCtrlZ = "\x1a"; Ascii_SUB = "\x1a" // substitute
                       Ascii_ESC = "\x1b" // escape
                       Ascii_FS  = "\x1c" // file separator
                       Ascii_GS  = "\x1d" // group separator
                       Ascii_RS  = "\x1e" // record separator
                       Ascii_US  = "\x1f" // unit separator
                       Ascii_DEL = "\x7f"
)

// ANSI input escape sequences emitted by keys.
// TODO: There are more edge cases, add them here!
const (
    KeyF1     = "\x1b\x5b\x59"
    KeyF2     = "\x1b\x5b\x60"
    KeyF3     = "\x1b\x5b\x61"
    KeyF4     = "\x1b\x5b\x62"
    KeyF5     = "\x1b\x5b\x31\x35\x7e"
    KeyF6     = "\x1b\x5b\x31\x37\x7e"
    KeyF7     = "\x1b\x5b\x31\x38\x7e"
    KeyF8     = "\x1b\x5b\x31\x39\x7e"
    KeyF9     = "\x1b\x5b\x32\x30\x7e"
    KeyF10    = "\x1b\x5b\x32\x31\x7e"
    KeyF11    = "\x1b\x5b\x32\x33\x7e"
    KeyF12    = "\x1b\x5b\x32\x34\x7e"

    KeyInsert = "\x1b\x5b\x32\x7e"
    KeyDelete = "\x1b\x5b\x33\x7e"
    KeyHome   = "\x1b\x5b\x48"
    KeyEnd    = "\x1b\x5b\x46"
    KeyPgUp   = "\x1b\x5b\x35\x7e"
    KeyPgDn   = "\x1b\x5b\x36\x7e"

    KeyUp     = "\x1b\x5b\x41"
    KeyDown   = "\x1b\x5b\x42"
    KeyRight  = "\x1b\x5b\x43"
    KeyLeft   = "\x1b\x5b\x44"
    KeyCtrlUp    = "\x1b\x5b\x31\x3b\x35\x41"
    KeyCtrlDown  = "\x1b\x5b\x31\x3b\x35\x42"
    KeyCtrlRight = "\x1b\x5b\x31\x3b\x35\x43"
    KeyCtrlLeft  = "\x1b\x5b\x31\x3b\x35\x44"
    KeyShiftUp    = "\x1b\x5b\x31\x3b\x32\x41"
    KeyShiftDown  = "\x1b\x5b\x31\x3b\x32\x42"
    KeyShiftRight = "\x1b\x5b\x31\x3b\x32\x43"
    KeyShiftLeft  = "\x1b\x5b\x31\x3b\x32\x44"
)

type Key struct {
    Value string
    Type int
}

// Stdin is buffered to allow logical read operations (e.g. ReadRune).
// Always read from terminal.Stdin instead of os.Stdin when also calling ReadKey().
// This is because terminal.Stdin may have additional bytes buffered.
var Stdin = bufio.NewReader(os.Stdin)

var keyAnsiLookup = map[string]bool{
    KeyF1: true, KeyF7:  true, KeyInsert: true,
    KeyF2: true, KeyF8:  true, KeyDelete: true,
    KeyF3: true, KeyF9:  true, KeyHome:   true,
    KeyF4: true, KeyF10: true, KeyEnd:    true,
    KeyF5: true, KeyF11: true, KeyPgUp:   true,
    KeyF6: true, KeyF12: true, KeyPgDn:   true,
    KeyUp:    true, KeyCtrlUp:    true, KeyShiftUp:    true,
    KeyDown:  true, KeyCtrlDown:  true, KeyShiftDown:  true,
    KeyLeft:  true, KeyCtrlRight: true, KeyShiftRight: true,
    KeyRight: true, KeyCtrlLeft:  true, KeyShiftLeft:  true,
}

// Blocks until Stdin can be read from, then returns the key read.
// A key represents (1) a Unicode rune or (2) an ANSI escape sequence.
// Note: An Alt modifier and an Escape followed by a rune are indistinguishable.
func ReadKey() Key {
    var key string
    for {
        key = ""
        for i := 0; i < 7; i++ {
            r, _, _ := Stdin.ReadRune()
            if r == unicode.ReplacementChar { break; }
            key += string(r)
            if key == "\x1b" && Stdin.Buffered() > 0 { continue; }
            if key == "\x1b\x5b" && Stdin.Buffered() > 0 { continue; }
            if keyAnsiLookup[key] { return Key{key, KEY_ANSI}; }
            if len(key) == 1 {
                if unicode.IsPrint(r) { return Key{key, KEY_RUNE}; }
                return Key{key, KEY_CTRL};
            }
            if len(key) == 2 && key[0] == 0x1b { return Key{key, KEY_ALT}; }
            if Stdin.Buffered() == 0 { break; }
        }
    }
}
