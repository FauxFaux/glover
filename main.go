package main

/*
extern void setup();
*/
import "C"

import (
    "bufio"
    "os/exec"
    "log"
    "os"
    "fmt"
)

func main() {
    cmd := exec.Command("keyhook/hook")
    stdout, err := cmd.StdoutPipe()
    if err != nil {
        log.Fatal(err)
    }
    C.setup()
    defer stdout.Close()
    err = cmd.Start()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("done!")
    scanner := bufio.NewScanner(stdout)
    for scanner.Scan() {
        fmt.Println(scanner.Text())
    }

    if err := scanner.Err(); err != nil {
        fmt.Fprintln(os.Stderr, "reading hook", err)
    }
}
