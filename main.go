package main

/*
#cgo CFLAGS: -DFOR_GO=1 -std=c99
extern int setup();
*/
import "C"

import (
	"bufio"
	"container/ring"
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

var toKeyName = map[StenoKey]string{
	LS:     "S",
	LT:     "T",
	LK:     "K",
	LP:     "P",
	LW:     "W",
	LH:     "H",
	LR:     "R",
	LA:     "A",
	LO:     "O",
	RE:     "E",
	RU:     "U",
	RF:     "F",
	RR:     "R",
	RP:     "P",
	RB:     "B",
	RL:     "L",
	RG:     "G",
	RT:     "T",
	RS:     "S",
	RD:     "D",
	RZ:     "Z",
	NUMBER: "#",
	STAR:   "*",
}

var numResolve = map[byte]StenoKey{
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
			log.Fatal("no "+name+".json: ", err, ", nor "+name+".json.template: ", err2)
		}
		fi = fi2
	}

	r := bufio.NewReader(fi)
	return json.NewDecoder(r), fi
}

func load(dict map[string]string) (chords Sequence) {
	chords = Sequence{Predecessors: map[Chord]*Sequence{}}
	loaded := 0
dicter:
	for k, v := range dict {
		splut := strings.Split(k, "/")
		var seq *Sequence = &chords
		for i := len(splut) - 1; i >= 0; i-- {
			ch, err := parseChord(splut[i])
			if nil != err {
				log.Println(err)
				continue dicter
			} else {
				var newVal *Sequence = seq.Predecessors[ch]
				if nil == newVal {
					newVal = &Sequence{Predecessors: map[Chord]*Sequence{}}
				}
				seq.Predecessors[ch] = newVal
				seq = newVal
			}
		}
		seq.Value = v
		loaded++
	}

	fmt.Println(loaded, "chords loaded")
	return
}

func render(ch Chord) string {
	s := ""
	for i := LS; i <= STAR; i++ {
		if 0 != (ch & h(i)) {
			if "" == s && i >= RF {
				s += "-"
			}
			s += toKeyName[i]
		}
	}
	return s
}

func lookup(chords Sequence, r *ring.Ring) (prop string) {
	r = r.Prev()

	// if there's no more inputs in the buffer,
	// and we got here, this must be what we want
	if nil == r.Value {
		return chords.Value
	}
	var ch Chord = r.Value.(Chord)

	// if we can go deeper, we should
	cand := chords.Predecessors[ch]
	if nil != cand {
		return lookup(*cand, r)
	}

	// if we can't go deeper, and we have something, it's right
	if "" != chords.Value {
		return chords.Value
	}

	return render(ch)
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
	var dict = map[string]string{}
	err = dictJ.Decode(&dict)
	if nil != err {
		log.Fatal("couldn't read dictionary: ", err)
	}

	chords := load(dict)

	keyMap := readKeys(c.Keys)

	for k, v := range chords.Predecessors {
		if v.Value == "have" {
			fmt.Printf("have %b\n", k)
		}
	}

	prev := ring.New(20)
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
		if 0 != c {
			prev.Value = c
			prev = prev.Next()
			fmt.Printf("%b: %s\n", c, lookup(chords, prev))
			if false {
				prev.Do(func(x interface{}) {
					fmt.Printf("%b, ", x)
				})
			}
		}
	}
}

/* vim: set noexpandtab: */
