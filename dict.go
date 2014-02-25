package main

import (
	"bufio"
	"container/ring"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
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

func lookup(chords Sequence, r *ring.Ring, depth int) (string, int) {
	r = r.Prev()

	// if there's no more inputs in the buffer,
	// and we got here, this must be what we want
	if nil == r.Value {
		return chords.Value, depth
	}
	var ch Chord = r.Value.(Chord)

	// if we can go deeper, we should
	cand := chords.Predecessors[ch]
	if nil != cand {
		return lookup(*cand, r, depth + 1)
	}

	// if we can't go deeper, and we have something, it's right
	if "" != chords.Value {
		return chords.Value, depth
	}

	return render(ch), 0
}

func readDict(name string) Sequence {
	dictJ, fi := openJson("dict")
	defer fi.Close()
	var dict = map[string]string{}
	err := dictJ.Decode(&dict)
	if nil != err {
		log.Fatal("couldn't read dictionary: ", err)
	}

	return load(dict)
}

/* vim: set noexpandtab: */
