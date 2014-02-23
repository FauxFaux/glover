package main

import (
    "testing"
)

func assertEquals(t *testing.T, expected, actual string) {
    if expected == actual {
        return
    }

    t.Fatalf("'%s' must equal '%s'", expected, actual)
}

func TestLoad(t *testing.T) {
    one := map[string]string {
        "SR": "have",
        "KWRES": "yes",
    }
    assertEquals(t, "yes", *load(one).Predecessors[130].Value)
}
