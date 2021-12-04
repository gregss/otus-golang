package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type wordf struct {
	word string
	frq  int
}

func Top10(s string) []string {
	var r []string

	if s == "" {
		return r
	}

	words := make(map[string]int)
	for _, word := range strings.Fields(s) {
		words[word]++
	}

	ws := make([]wordf, 0, len(words))
	for word, frq := range words {
		ws = append(ws, wordf{word: word, frq: frq})
	}

	sort.Slice(ws, func(i, j int) bool { return ws[i].frq > ws[j].frq })

	var lexs []wordf
	pfrq := ws[0].frq
	for i, w := range ws {
		if pfrq != w.frq {
			r = applyLexSlice(lexs, r)

			lexs = nil
			lexs = append(lexs, w)
			pfrq = w.frq
		} else {
			lexs = append(lexs, w)
		}

		if len(r) > 10 {
			r = r[:10]
			break
		}

		if i == len(ws)-1 {
			r = applyLexSlice(lexs, r)
		}
	}

	return r
}

func applyLexSlice(lexs []wordf, r []string) []string {
	sort.Slice(lexs, func(i, j int) bool {
		return lexs[i].word < lexs[j].word
	})
	for _, w := range lexs {
		r = append(r, w.word)
	}

	return r
}
