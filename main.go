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
	"strings"
)

var output = make(chan string, 100)

//export chordReleased
func chordReleased(bytes *C.char) {
	output <- C.GoString(bytes)
}

const sequenceLength = 10

type StenoKey uint8
type QwertyKey uint8
type VKey uint8
type Chord uint32
type Sequence [sequenceLength]Chord

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

var toEnum = map[string]StenoKey {
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

var numResolve = map[byte]StenoKey {
	'1': LS,
	'2': LT,
	'3': LP,
	'4': LH,
	'5': LA,
	'6': RF,
	'7': RP,
	'8': RL,
	'9': RT,
	'0': RD,
	'#': NUMBER,
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

func openJson(name string) (dec *json.Decoder, fi *os.File) {
	fi, err := os.Open(name + ".json")
	if nil != err {
		fi2, err2 := os.Open(name + ".json.template")
		if nil != err2 {
			log.Fatal("no " + name + ".json: ", err, ", nor " + name + ".json.template: ", err2)
		}
		fi = fi2
	}

	r := bufio.NewReader(fi)
	return json.NewDecoder(r), fi
}

func h(s StenoKey) Chord {
	return 1 << s;
}

func splitChar(c byte, passedMid bool) (ret Chord, newMid bool, err error) {
	switch c {
		case 'A':
			return h(LA), true, nil
		case 'O':
			return h(LO), true, nil
		case 'E':
			return h(RE), true, nil
		case 'U':
			return h(RU), true, nil
		case '*':
			return h(STAR), true, nil
		case '-':
			return 0, true, nil
	}

	foundNum := numResolve[c]
	if 0 != foundNum {
		return h(NUMBER) | h(foundNum), true, nil
	}

	var encoding string
	if passedMid {
		encoding = fmt.Sprintf("-%c", c)
	} else {
		encoding = fmt.Sprintf("%c-", c)
	}

	found := toEnum[encoding]
	if 0 != found {
		return h(found), passedMid, nil
	}

	return 0, passedMid, fmt.Errorf("not happy with '%c' %t", c, passedMid)
}

func parseChord(s string) (ret Chord, err error) {
	passedMid := false
	for i := 0; i < len(s); i++ {
		ch, newMid, err := splitChar(s[i], passedMid)
		passedMid = newMid
		if nil != err {
			return 0, fmt.Errorf("Can't read '%c' in '%s': %s", s[i], s, err)
		}
		ret |= ch
	}
	return ret, nil
}

func main() {
	go C.setup()
	dec, fi := openJson("config")
	defer fi.Close()

	type Config struct {
		Keys map[string]string
	}

	var c Config
	err := dec.Decode(&c)
	if nil != err {
		log.Fatal("couldn't unmarshal: ", err)
	}

	dictJ, fi := openJson("dict")
	defer fi.Close()
	var dict = map[string]string {}
	err = dictJ.Decode(&dict)
	if nil != err {
		log.Fatal("couldn't read dictionary: ", err)
	}

	var chords = map[Sequence]string {}
dicter:
	for k, v := range dict {
		splut := strings.Split(k, "/")
		if len(splut) > sequenceLength {
			log.Println("can't cope with", k, ":", v, "as it's too long")
			continue
		}
		chs := Sequence {}
		for i, part := range splut {
			ch, err := parseChord(part)
			if nil != err {
				log.Println(err)
				continue dicter
			} else {
				chs[i] = ch
			}
		}
		chords[chs] = v
	}

	fmt.Println(len(chords), "chords loaded")

	var keyMap = map[VKey]StenoKey{}
	for key, value := range c.Keys {
		var val VKey = keyNameToVKey(key)
		var x StenoKey = toEnum[value]
		if 0 == x {
			log.Fatalf("I can't parse '%s' as a steno key", value)
		}

		keyMap[val] = x
	}

	for k, v := range chords {
		if v == "have" {
			fmt.Printf("have %b\n", k)
		}
	}

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
//			fmt.Printf("%c: %s;   ", vk, fromEnum[sk])
		}
		if 0 != c {
			fmt.Printf("%b: %s\n", c, chords[Sequence{c}])
		}
	}
}

/* vim: set noexpandtab: */
