package common

import (
	"unicode"
	"unicode/utf8"
)

func Match(source, target string) bool {
	lenDiff := len(target) - len(source)

	if lenDiff < 0 {
		return false
	}

	if lenDiff == 0 && source == target {
		return true
	}

Outer:
	for _, r1 := range source {
		for i, r2 := range target {
			if unicode.ToLower(r1) == unicode.ToLower(r2) {
				target = target[i+utf8.RuneLen(r2):]
				continue Outer
			}
		}
		return false
	}

	return true
}
