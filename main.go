package main

/*
#cgo CFLAGS: -DFOR_GO=1 -std=c99
extern int setup();
*/
import "C"

import (
    "fmt"
    "bufio"
    "log"
    "os"
    "strconv"
    "encoding/json"
)

var output = make(chan string, 100)

//export chordReleased
func chordReleased(bytes *C.char) {
    output <- C.GoString(bytes)
}

const (
    LS = iota
    LT
    LK
    LP
    LW
    LH
    LR
    A
    O
    E
    U
    RF
    RR
    RP
    RB
    RL
    RG
    RT
    RS
    RD
    RZ
)

func main() {
    go C.setup();
    fi, err := os.Open("config.json")
    if nil != err {
        fi2, err2 := os.Open("config.json.template")
        if nil != err2 {
            log.Fatal("no config.json: ", err, ", nor config.json.template: ", err2)
        }
        fi = fi2
    }

    r := bufio.NewReader(fi)
    dec := json.NewDecoder(r)

    type Config struct {
        Keys map[string]string
    }

    var c Config
    err = dec.Decode(&c)
    if nil != err {
        log.Fatal("couldn't unmarshal: ", err)
    }

    var keyMap [255]string
    for key, value := range c.Keys {
        val, err := strconv.Atoi(key)
        if nil != err {
            log.Fatal("vkeys must be integers: ", err)
        }
        keyMap[val] = value
    }

    for {
        req := <-output
        // this cannot use range as it is not utf-8
        for i := 0; i < len(req); i++ {
            fmt.Printf("%d: %s ", req[i], keyMap[req[i]])
        }
        fmt.Println()
    }
}

