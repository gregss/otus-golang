package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode"
)

type symbol struct {
	r     rune
	multi bool
	esc   bool
}

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	var b strings.Builder
	var crnt, prev symbol
	for i, r := range s {
		crnt = symbol{
			r:     r,
			multi: unicode.IsDigit(r) && !prev.esc,
			esc:   string(r) == `\` && !prev.esc,
		}

		if crnt.multi {
			if (symbol{} == prev || prev.multi) {
				return "", ErrInvalidString
			}
			c := int(crnt.r - '0')
			if c > 0 {
				b.WriteString(strings.Repeat(string(prev.r), c))
			}
		} else if !prev.multi && !prev.esc && i != 0 {
			b.WriteRune(prev.r)
		}

		prev = crnt
	}

	if !crnt.multi && !crnt.esc && (symbol{} != prev) {
		b.WriteRune(crnt.r)
	}

	return b.String(), nil
}
