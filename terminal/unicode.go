package terminal

// TODO: Better Unicode support.
// Char is a list of runes.
// Provide utilities to read chars, not runes.
// Implement important string utilities: iter, length, etc.

// TODO: Unicode is difficult.
// Byte > Rune > Glyph (= Char) > String.
// Usually 1 rune = 1 glyph, but not always (e.g. emojis)!
// For example len() returns the number of bytes in the string.
// Storing slices of runes would return the number of runes.
// https://github.com/rivo/uniseg
// Unfortunately Go makes this process very tedious...

// TODO: The editor should use slices of chars!
