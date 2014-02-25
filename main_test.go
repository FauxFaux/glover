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

func assertIEquals(t *testing.T, expected, actual int) {
	if expected == actual {
		return
	}

	t.Fatalf("'%d' must equal '%d'", actual, expected)
}

func dump(s Sequence, depth int) {
	fmt.Println(s.Value)
	for k, v := range s.Predecessors {
		fmt.Printf("%d %9b:\n", depth, k)
		dump(*v, depth+1)
	}
}

var one = map[string]string{
	"SR":            "have",
	"KWRES":         "yes",
	"SEP/RAT":       "separate",
	"RAT":           "rat",
	"S/S/S/S/S/S/S": "Jodie Foster",
}

const S = Chord(2)
const T = Chord(4)
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

func assertKill(t *testing.T, expected int, s Sequence, r *ring.Ring) string {
	ret, kill := lookup(s, r, -1)
	assertIEquals(t, expected, kill)
	return ret
}

func TestLookup(t *testing.T) {
	out := load(one)
	r := ring.New(5)
	r.Value = SR
	r = r.Next()
	assertEquals(t, "have", assertKill(t, 0, out, r))
	r.Value = RAT
	r = r.Next()
	assertEquals(t, "rat", assertKill(t, 0, out, r))
	r.Value = SEP
	r = r.Next()
	assertEquals(t, "SEP", assertKill(t, 0, out, r))
	r.Value = RAT
	r = r.Next()
	assertEquals(t, "separate", assertKill(t, 1, out, r))
}

func TestMistranslate(t *testing.T) {
	out := load(one)
	r := ring.New(20)
	for i := 0; i < 10; i++ {
		r.Value = SR; r = r.Next()
	}

	r.Value = T; r = r.Next()

	assertEquals(t, "T", assertKill(t, 0, out, r))
}

func TestJodieFoster(t *testing.T) {
	out := load(one)
	r := ring.New(20)
	for i := 0; i < 7; i++ {
		r.Value = S
		r = r.Next()
	}
	assertEquals(t, "Jodie Foster", assertKill(t, 6, out, r))
}

/* vim: set noexpandtab: */
