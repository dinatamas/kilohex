// +build !darwin,!freebsd,!netbsd,!openbsd,!windows

package terminal

import (
    "golang.org/x/sys/unix"
)

const (
    TIOGET = unix.TCGETS
    TIOSET = unix.TCSETS
)
