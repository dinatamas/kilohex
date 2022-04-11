KiloGo
======

  - My version of the wonderful Kilo text editor in Go.

Backlog
-------

  - Handle tab characters properly.

  - Refactor to be event-driven and use polling.

  - React to SIGWINCH events and do not use DetectResize().

  - Basic terminfo checking for supported escape sequences?

References
----------

  [1] https://github.com/antirez/kilo
  [2] https://github.com/srinathh/gokilo
  [3] https://viewsourcecode.org/snaptoken/kilo/
  [4] https://github.com/nsf/termbox-go
  [5] https://pkg.go.dev/golang.org/x/term
  [6] https://gist.github.com/fnky/458719343aabd01cfb17a3a4f7296797
  [7] https://invisible-island.net/xterm/ctlseqs/ctlseqs.html
