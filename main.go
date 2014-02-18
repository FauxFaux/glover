package main

/*
#cgo CPPFLAGS: -DNO_MAIN=1
extern int setup();
*/
import "C"

import (
)

func main() {
    go C.setup();
}
