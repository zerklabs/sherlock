package sherlock

import (
	"strings"
)

var (
	alphabet = map[string]int{
		"a": 1,
		"A": 1,
		"b": 2,
		"B": 2,
		"c": 3,
		"C": 3,
		"d": 4,
		"D": 4,
		"e": 5,
		"E": 5,
		"f": 6,
		"F": 6,
		"g": 7,
		"G": 7,
		"h": 8,
		"H": 8,
		"i": 9,
		"I": 9,
		"j": 10,
		"J": 10,
		"k": 11,
		"K": 11,
		"l": 12,
		"L": 12,
		"m": 13,
		"M": 13,
		"n": 14,
		"N": 14,
		"o": 15,
		"O": 15,
		"p": 16,
		"P": 16,
		"q": 17,
		"Q": 17,
		"r": 18,
		"R": 18,
		"s": 19,
		"S": 19,
		"t": 20,
		"T": 20,
		"u": 21,
		"U": 21,
		"v": 22,
		"V": 22,
		"w": 23,
		"W": 23,
		"x": 24,
		"X": 24,
		"y": 25,
		"Y": 25,
		"z": 26,
		"Z": 26,
	}
)

func TokenizeLine(line string) []string {
	return strings.Split(line, " ")
}

func ScoreWord(word string) float32 {
	var score float32
	// wordLength := len(word)

	score = 0.0

	split := strings.Split(word, "")

	// for keeping track of the last position in the
	// alphabet for the last character seen
	lastPosition := 0

	for _, v := range split {
		// currently not ranking non-alphabetic characters
		if alphabet[v] > 0 {
			if lastPosition == 0 {
				lastPosition = alphabet[v]
			} else {
				dist := lastPosition - alphabet[v]

				if dist < 0 {
					dist = -dist
				}

				percentDist := (dist / 26) * 100

				if percentDist >= 50 {
					score = score - 0.1
				} else {
					score = score + 0.2
				}
			}
		} else {
			// partial loss for non-captured characters
			score = score - 0.1
		}
	}

	return score
}

func ClassifyScore(score float32) int {
	if score <= 0.3 {
		return -1
	}

	if score <= 0.6 {
		return 0
	}

	return 1
}