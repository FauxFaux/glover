package main

/*
#cgo CFLAGS: -DFOR_GO=1 -std=c99
extern int setup();
*/
import "C"

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

var output = make(chan string, 100)

//export chordReleased
func chordReleased(bytes *C.char) {
	output <- C.GoString(bytes)
}

type StenoKey uint8
type QwertyKey uint8
type VKey uint8

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

var toEnum = map[string]StenoKey{
	"S-": LS,
	"T-": LT,
	"K-": LK,
	"P-": LP,
	"W-": LW,
	"H-": LH,
	"R-": LR,
	"A-": LA,
	"O-": LO,
	"-E": RE,
	"-U": RU,
	"-F": RF,
	"-R": RR,
	"-P": RP,
	"-B": RB,
	"-L": RL,
	"-G": RG,
	"-T": RT,
	"-S": RS,
	"-D": RD,
	"-Z": RZ,
	"#":  NUMBER,
	"*":  STAR,
}

func reverse(in map[string]StenoKey) map[StenoKey]string {
	ret := map[StenoKey]string{}
	for key, value := range in {
		ret[value] = key
	}
	return ret
}

var fromEnum = reverse(toEnum)

func vkeyToChar(in VKey) QwertyKey {
	if in >= 65 && in < 91 {
		return QwertyKey(in - 65 + 'a')
	}

	if in >= 48 && in < 58 {
		return QwertyKey(in - 48 + '0')
	}

	return 0
}

func keyNameToVKey(in string) VKey {
	var from QwertyKey = QwertyKey(in[0])
	if from >= 'a' && from <= 'z' {
		return VKey(from - 'a' + 65)
	}

	if from >= '0' && from <= '9' {
		return VKey(from - '0' + 0x30)
	}

	switch from {
	case ';':
		return 0xba
	case '[':
		return 0xdb
	case ']':
		return 0xdd
	case '-':
		return 0xbd
	case '\'':
		return 0xde
	}

	log.Fatalf("don't understand '%c' as a VKey", from)
	return 0
}

func openJson(name string) *json.Decoder {
	fi, err := os.Open(name + ".json")
	if nil != err {
		fi2, err2 := os.Open(name + ".json.template")
		if nil != err2 {
			log.Fatal("no " + name + ".json: ", err, ", nor " + name + ".json.template: ", err2)
		}
		fi = fi2
	}

	defer fi.Close()

	r := bufio.NewReader(fi)
	return json.NewDecoder(r)
}

func main() {
	go C.setup()
	dec := openJson("config")

	type Config struct {
		Keys map[string]string
	}

	var c Config
	err := dec.Decode(&c)
	if nil != err {
		log.Fatal("couldn't unmarshal: ", err)
	}

	var keyMap = map[VKey]StenoKey{}
	for key, value := range c.Keys {
		var val VKey = keyNameToVKey(key)
		var x StenoKey = toEnum[value]
		if 0 == x {
			log.Fatalf("I can't parse '%s' as a steno key", value)
		}

		keyMap[val] = x
	}

	for {
		req := <-output
		// this cannot use range as it is not utf-8
		for i := 0; i < len(req); i++ {
			var vk VKey = VKey(req[i])
			var sk StenoKey = keyMap[vk]
			if 0 == sk {
				continue
			}
			fmt.Printf("%c: %s;   ", vk, fromEnum[sk])
		}
		fmt.Println()
	}
}

/* vim: set noexpandtab: */
