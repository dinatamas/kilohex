package terminal

// ANSI output escape sequences.
const (
    AnsiClearLine   = "\x1b\x5b\x4b"
    AnsiClearScreen = "\x1b\x5b\x32\x4a"
    AnsiClearBelow  = "\x1b\x5b\x30\x4a"

    // Cursor saving and restoring: https://stackoverflow.com/a/29163244
    AnsiCursorGetPosition = "\x1b\x5b\x36\x6e"
    AnsiCursorHide        = "\x1b\x5b\x3f\x32\x35\x31"
    AnsiCursorHome        = "\x1b\x5b\x48"
    AnsiCursorRestore     = "\x1b\x38"
    AnsiCursorSave        = "\x1b\x37"
    AnsiCursorShow        = "\x1b\x5b\x3f\x32\x35\x68"

    AnsiInverse      = "\x1b\x5b\x37\x6d"
    AnsiInverseReset = "\x1b\x5b\x32\x37\x6d"
)

// ANSI output escape sequence templates.
const (
    AnsiCursorSetPosition = "\x1b\x5b%d\x3b%d\x48"
)

// ANSI input escape sequence templates.
const (
    AnsiCursorPositionReport = "\x1b\x5b%d\x3b%d\x52"
)
