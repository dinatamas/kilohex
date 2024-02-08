KiloHex
=======

Modified version of the Kilo text editor - as a hex editor - written in Go!

Tasklist
--------

**Short term:**

- [ ] Preserve original file permissions when writing.
- [ ] Add whitespace between bytes.
- [ ] Display ASCII representation on the right side.
- [ ] Display row offset in hex on the left side.
- [ ] Display current byte position in the current row.
- [ ] Color different byte classes.
- [ ] Allow only hex character changes.

**Medium term:**

- [ ] Display column offsets in hex in the top row.
- [ ] Support byte deletion.
- [ ] Support byte insertion.

**Long term:**

- [ ] Introduce binary display modes.
- [ ] Support bit operations.

References
----------

- [kilo](            https://github.com/antirez/kilo)
- [gokilo](          https://github.com/srinathh/gokilo)
- [kilo explained](  https://viewsourcecode.org/snaptoken/kilo/)
- [termbox-go](      https://github.com/nsf/termbox-go)
- [go term](         https://pkg.go.dev/golang.org/x/term)
- [docker term](     https://github.com/moby/term/)
- [ansi escape](     https://gist.github.com/fnky/458719343aabd01cfb17a3a4f7296797)
- [xterm ctlseqs](   https://invisible-island.net/xterm/ctlseqs/ctlseqs.html)
- [hexyl](           https://github.com/sharkdp/hexyl)
