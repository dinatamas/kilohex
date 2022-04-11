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

    if filename == "" {
        log.Println("usage: kilogo <filename>")
        return 1
    }

    // Read current file content.
    bytebuf, err := ioutil.ReadFile(filename)
    if errors.Is(err, os.ErrNotExist) {
        bytebuf = []byte{}
    } else if err != nil {
        log.Print("error: ", err)
        return 2
    }
    buffer := string(bytebuf)

    // Run the editor.
    editor := NewEditor(filename, buffer)
    defer editor.Restore()
    buffer = editor.Run()

    // Write the edited content to the file.
    err = ioutil.WriteFile(filename, []byte(buffer), 0666)
    if err != nil {
        log.Print("error: ", err)
        return 4
    }

    return 0
}
