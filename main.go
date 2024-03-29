package main

/*
#cgo CFLAGS: -DFOR_GO=1 -std=c99
extern int setup();
*/
import "C"

import (
	"container/ring"
	"fmt"
)

var output = make(chan string, 100)

//export chordReleased
func chordReleased(bytes *C.char) {
	output <- C.GoString(bytes)
}

type StenoKey uint8
type QwertyKey uint8
type VKey uint8
type Chord uint32

type Sequence struct {
	Value        string
	Predecessors map[Chord]*Sequence
}

const (
	_ StenoKey = iota
	LS
	LT
	LK
	LP
	LW
	LH
	LR
	LA
	LO
	RE
	RU
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
	NUMBER
	STAR
)

func vkeyToChar(in VKey) QwertyKey {
	if in >= 65 && in < 91 {
		return QwertyKey(in - 65 + 'a')
	}

	if in >= 48 && in < 58 {
		return QwertyKey(in - 48 + '0')
	}

	return 0
}

func main() {
	go C.setup()

	keyMap := readKeyMap("config")
	chords := readDict("dict")

	prevChords := ring.New(20)
	prevOutput := ring.New(20)
	for {
		req := <-output
		// this cannot use range as it is not utf-8
		var c Chord
		for i := 0; i < len(req); i++ {
			var vk VKey = VKey(req[i])
			var sk StenoKey = keyMap[vk]
			if 0 == sk {
				continue
			}
			c |= h(sk)
		}
		if 0 == c {
			continue
		}
		prevChords.Value = c
		prevChords = prevChords.Next()
		res, kill := lookup(chords, prevChords, -1)
		for i := 0; i < kill; i++ {
			prevOutput = prevOutput.Prev()
			fmt.Printf("kill %s\n", prevOutput.Value)
		}
		prevOutput.Value = res
		prevOutput = prevOutput.Next()
		fmt.Printf("%s: %s (%d)\n", render(c), res, kill)
	}
}

/* vim: set noexpandtab: */
