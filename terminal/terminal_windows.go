// +build windows

package terminal

import (
    w "golang.org/x/sys/windows"
)

type TerminalConfig struct {
    InputFlags  uint32
    OutputFlags uint32
}

func GetConfig() TerminalConfig {
    config := TerminalConfig{}
    w.GetConsoleMode(w.Stdin, &config.InputFlags)
    w.GetConsoleMode(w.Stdout, &config.OutputFlags)
    return config
}

func GetSize() (int, int) {
    info := w.ConsoleScreenBufferInfo{}
    w.GetConsoleScreenBufferInfo(w.Stdout, &info)
    return info.Window.Bottom - info.Window.Top + 1,
           info.Windor.Right - info.Window.Left + 1
}

func SetConfig(config TerminalConfig) {
    w.SetConsoleMode(w.Stdin, config.InputFlags)
    w.SetConsoleMode(w.Stdout, config.OutputFlags)
}

func SetRawMode() {
    w.SetConsoleMode(
        w.Stdin,
        w.ENABLE_EXTENDED_FLAGS |
        w.ENABLE_VIRTUAL_TERMINAL_INPUT,
    )
    w.SetConsoleMode(
        w.Stdout,
        w.ENABLE_VIRTUAL_TERMINAL_PROCESSING |
        w.ENABLE_PROCESSED_OUTPUT |
        w.DISABLE_NEWLINE_AUTO_RETURN,
    )
}
