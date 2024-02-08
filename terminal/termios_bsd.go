// +build darwin freebsd openbsd netbsd

package terminal

import (
    "golang.org/x/sys/unix"
)

const (
    TIOGET = unix.TIOCGETA
    TIOSET = unix.TIOCSETA
)
