package main

import (
	"fmt"
	"log"
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

func h(s StenoKey) Chord {
	return 1 << s
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

func readKeys(keys map[string]string) map[VKey]StenoKey {
	var keyMap = map[VKey]StenoKey{}
	for key, value := range keys {
		var val VKey = keyNameToVKey(key)
		var x StenoKey = toEnum[value]
		if 0 == x {
			log.Fatalf("I can't parse '%s' as a steno key", value)
		}

		keyMap[val] = x
	}
	return keyMap
}

func readKeyMap(name string) map[VKey]StenoKey {
	dec, fi := openJson(name)
	defer fi.Close()

	type Config struct {
		Keys map[string]string
	}

	var c Config
	err := dec.Decode(&c)
	if nil != err {
		log.Fatal("couldn't unmarshal: ", err)
	}

	return readKeys(c.Keys)
}

/* vim: set noexpandtab: */
