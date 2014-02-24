package main

import (
	"container/ring"
	"fmt"
	"testing"
)

func assertEquals(t *testing.T, expected, actual string) {
	if expected == actual {
		return
	}

	t.Fatalf("'%s' must equal '%s'", actual, expected)
}

func dump(s Sequence, depth int) {
	fmt.Println(s.Value)
	for k, v := range s.Predecessors {
		fmt.Printf("%d %9b:\n", depth, k)
		dump(*v, depth+1)
	}
}

var one = map[string]string{
	"SR":      "have",
	"KWRES":   "yes",
	"SEP/RAT": "separate",
	"RAT":     "rat",
}

const SR = Chord(128 + 2)
const KWRES = Chord(525480)
const RAT = Chord(262528)
const SEP = Chord(17410)

func TestLoad(t *testing.T) {
	out := load(one)

	dump(out, 0)

	assertEquals(t, "have", out.Predecessors[SR].Value)
	assertEquals(t, "yes", out.Predecessors[KWRES].Value)
	assertEquals(t, "separate", out.Predecessors[RAT].Predecessors[SEP].Value)
	assertEquals(t, "rat", out.Predecessors[RAT].Value)
}

func TestLookup(t *testing.T) {
	out := load(one)
	r := ring.New(5)
	r.Value = SR
	r = r.Next()
	assertEquals(t, "have", lookup(out, r))
	r.Value = RAT
	r = r.Next()
	assertEquals(t, "rat", lookup(out, r))
	r.Value = SEP
	r = r.Next()
	assertEquals(t, "SEP", lookup(out, r))
	r.Value = RAT
	r = r.Next()
	assertEquals(t, "separate", lookup(out, r))
}

/* vim: set noexpandtab: */
