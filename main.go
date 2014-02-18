package main

/*
#cgo CFLAGS: -DFOR_GO=1 -std=c99
extern int setup();
*/
import "C"

import (
    "fmt"
)

var output = make(chan string, 100)

//export chordReleased
func chordReleased(bytes *C.char) {
    output <- C.GoString(bytes)
}

func main() {
    go C.setup();
    for {
        req := <-output
        // this cannot use range as it is not utf-8
        for i := 0; i < len(req); i++ {
            fmt.Printf("%x ", req[i])
        }
        fmt.Println()
    }
}

