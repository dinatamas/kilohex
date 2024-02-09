KiloHex
=======

Modified version of the Kilo text editor - as a hex editor - written in Go!

Tasklist
--------

**Missing:**

- [ ] Double cursor on bytes.
- [ ] Implement byte editing.
- [ ] Reintroduce file saving.

**Short term:**

- [ ] Extra cursor on decoded content.
- [ ] Preserve original file permissions when writing.
- [ ] Display current byte position in the current row.
- [ ] Color different byte classes.
- [ ] Utils: `SplitInto()`

**Medium term:**

- [ ] Display column offsets in hex in the top row.
- [ ] Support byte deletion.
- [ ] Support byte insertion.

**Long term:**

- [ ] Introduce hexdump command line arguments.
- [ ] Introduce binary display modes.
- [ ] Support bit operations.
- [ ] Full Unicode support (graphemes instead of runes).

**Bugs:**

- [ ] After resizing the left X-offset is not corrected.

References
----------

- [kilo](                https://github.com/antirez/kilo)
- [gokilo](              https://github.com/srinathh/gokilo)
- [kilo explained](      https://viewsourcecode.org/snaptoken/kilo/)
- [termbox-go](          https://github.com/nsf/termbox-go)
- [go term](             https://pkg.go.dev/golang.org/x/term)
- [docker term](         https://github.com/moby/term/)
- [ansi escape](         https://gist.github.com/fnky/458719343aabd01cfb17a3a4f7296797)
- [xterm ctlseqs](       https://invisible-island.net/xterm/ctlseqs/ctlseqs.html)
- [hexyl](               https://github.com/sharkdp/hexyl)
- [text normalization](  https://go.dev/blog/normalization)
- [Unicode problems](    https://stackoverflow.com/a/12668840)
- [uniseg](              https://github.com/rivo/uniseg)
