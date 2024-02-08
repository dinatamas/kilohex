package main

import (
    "errors"
    "flag"
    "io/ioutil"
    "log"
    "os"
)

func main() {
    os.Exit(_main())
}

func _main() int {
    // Parse command line arguments.
    flag.Parse()
    filename := flag.Arg(0)

    log.SetFlags(0)
    if filename == "" {
        log.Println("usage: kilohex <filename>")
        return 1
    }

    // Read current file content.
    buffer, err := ioutil.ReadFile(filename)
    if errors.Is(err, os.ErrNotExist) {
        buffer = []byte{}
    } else if err != nil {
        log.Print("file read error: ", err)
        return 2
    }

    // Run the editor.
    editor := NewEditor(filename, buffer)
    defer editor.Restore()
    buffer, err = editor.Run()
    if err != nil {
        log.Print("editor error: ", err)
        return 3
    }

    // Write the edited content to the file.
    err = ioutil.WriteFile(filename, buffer, 0666)
    if err != nil {
        log.Print("file write error: ", err)
        return 4
    }

    return 0
}
