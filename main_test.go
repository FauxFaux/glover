package main

import (
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

func TestLoad(t *testing.T) {
	one := map[string]string{
		"SR":      "have",
		"KWRES":   "yes",
		"SEP/RAT": "separate",
		"RAT":     "rat",
	}

	out := load(one)

	dump(out, 0)

	SR := Chord(128 + 2)
	KWRES := Chord(525480)
	RAT := Chord(262528)
	SEP := Chord(17410)

	assertEquals(t, "have", out.Predecessors[SR].Value)
	assertEquals(t, "yes", out.Predecessors[KWRES].Value)
	assertEquals(t, "separate", out.Predecessors[RAT].Predecessors[SEP].Value)
	assertEquals(t, "rat", out.Predecessors[RAT].Value)
}
