// +build !windows

package terminal

import (
    x "golang.org/x/sys/unix"
)

type TerminalConfig *x.Termios

func GetConfig() TerminalConfig {
    termios, _ := x.IoctlGetTermios(x.Stdin, x.TCGETS)
    return termios
}

func GetSize() (int, int) {
    winSize, _ := x.IoctlGetWinsize(x.Stdout, x.TIOCGWINSZ)
    return int(winSize.Row), int(winSize.Col)
}

func SetConfig(config TerminalConfig) {
    x.IoctlSetTermios(x.Stdin, x.TCSETSF, config)
}

func SetRawMode() {
    termios, _ := x.IoctlGetTermios(x.Stdin, x.TCGETS)
    termios.Lflag = termios.Lflag &^ (x.ECHO | x.ICANON | x.ISIG | x.IEXTEN)
    termios.Iflag = termios.Iflag &^ (x.IXON | x.ICRNL | x.BRKINT | x.INPCK | x.ISTRIP)
    termios.Oflag = termios.Oflag &^ x.OPOST
    termios.Cflag = termios.Cflag | x.CS8
    x.IoctlSetTermios(x.Stdin, x.TCSETSF, termios)
}
