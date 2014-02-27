package sherlock

import (
	"bytes"
	"encoding/gob"
	"log"
	"math"
	"strings"
)

var (
	alphabet = map[string]int{
		"a":  1,
		"A":  1,
		"b":  2,
		"B":  2,
		"c":  3,
		"C":  3,
		"d":  4,
		"D":  4,
		"e":  5,
		"E":  5,
		"f":  6,
		"F":  6,
		"g":  7,
		"G":  7,
		"h":  8,
		"H":  8,
		"i":  9,
		"I":  9,
		"j":  10,
		"J":  10,
		"k":  11,
		"K":  11,
		"l":  12,
		"L":  12,
		"m":  13,
		"M":  13,
		"n":  14,
		"N":  14,
		"o":  15,
		"O":  15,
		"p":  16,
		"P":  16,
		"q":  17,
		"Q":  17,
		"r":  18,
		"R":  18,
		"s":  19,
		"S":  19,
		"t":  20,
		"T":  20,
		"u":  21,
		"U":  21,
		"v":  22,
		"V":  22,
		"w":  23,
		"W":  23,
		"x":  24,
		"X":  24,
		"y":  25,
		"Y":  25,
		"z":  26,
		"Z":  26,
		"0":  27,
		"1":  28,
		"2":  29,
		"3":  30,
		"4":  31,
		"5":  32,
		"6":  33,
		"7":  34,
		"8":  35,
		"9":  36,
		":":  37,
		"-":  38,
		"\\": 39,
		"/":  40,
		"$":  41,
		".":  42,
		",":  43,
		"!":  44,
		"@":  45,
		"#":  46,
		"%":  47,
		"^":  48,
		"&":  49,
		"*":  50,
		"(":  51,
		")":  52,
		"_":  53,
		"+":  54,
		"=":  55,
		"<":  56,
		">":  57,
		"?":  58,
		"\"": 59,
		"'":  60,
		";":  61,
		"|":  62,
		"[":  63,
		"]":  64,
		"{":  65,
		"}":  66,
	}
)

type TDIDF struct {
	Word                     string
	RawTermFrequency         int
	LineFrequency            int
	TermFrequency            float64
	InverseDocumentFrequency float64
	Total                    float64
	Score                    float64
	Classification           int
	LeftNeighbor             *TDIDF
	RightNeighbor            *TDIDF
}

func TokenizeLine(line string) []string {
	return strings.Split(line, " ")
}

// Calculate the raw frequency (number of occurrences) of a given word
func (self *TDIDF) RawFrequency(doc *string) {
	self.RawTermFrequency = strings.Count(*doc, self.Word)
}

// Calculate the TD-IDF
func (self *TDIDF) Frequency(maxRawFrequency int, totalWords int) {
	// augmented frequency to prevent bias towards longer documents (# of lines)
	self.TermFrequency = float64(float64(0.5) + ((float64(0.5) * float64(self.RawTermFrequency)) / float64(maxRawFrequency)))

	// using the logarithmically scaled frequency
	self.InverseDocumentFrequency = math.Log(float64(totalWords / (1.0 + self.RawTermFrequency)))
	self.Total = (self.TermFrequency * self.InverseDocumentFrequency)
}

// Using a variation of [k-NN](http://en.wikipedia.org/wiki/K-nearest_neighbors_algorithm),
// rank words based on the distances in the known alphabet of each of its characters
func (self *TDIDF) ScoreWord() {
	var score float64
	var unresolvedLoss float64
	var unresolvedPercent int
	var unresolved int
	var percentDist int
	var lastPosition int
	var wordLength int

	// store the length of the given word
	wordLength = len(self.Word)
	// for keeping track of the last position in the
	// alphabet for the last character seen
	lastPosition = 0

	// starting score
	score = 0.0

	// starting unresolved loss
	unresolvedLoss = 0.0

	// start characters that were unresolved
	unresolved = 0

	if wordLength == 0 {
		self.Score = 0.0
	} else if wordLength < 4 {
		self.Score = 0.0
	} else {

		// break the word up into characters
		split := strings.Split(self.Word, "")

		for _, v := range split {
			// currently not ranking non-alphabetic characters
			if alphabet[v] > 0 {
				if lastPosition == 0 {
					lastPosition = alphabet[v]
				} else {
					dist := lastPosition - alphabet[v]

					// get the absolute value of dist
					if dist < 0 {
						dist = -dist
					}

					percentDist = (dist / len(alphabet)) * 100

					// if the current letter and the last are further than 50% away,
					// calculate a partial loss
					if percentDist >= 50 {
						score = score - 0.2
					} else if percentDist >= 0 {
						// otherwise, improve the score of this word
						score = score + 0.2
					} else if percentDist == 0 {
						score = score - 0.4
					}
				}
			} else {
				// partial loss for non-captured characters
				score = score - 0.1

				unresolved = unresolved + 1
			}

			// set the last position to this character
			lastPosition = alphabet[v]
		}

		unresolvedPercent = (unresolved / wordLength) * 100

		for i := 0; i < unresolved; i++ {
			unresolvedLoss += 0.1
		}

		// subtract the unresolved loss from the overall score
		// 0.1 penalty per unresolved character
		score = score - unresolvedLoss

		// penalty for a larger set of unresolved characters
		if unresolvedPercent >= 50 {
			score = score - 0.2
		}

		self.Score = score
	}
}

func (self *TDIDF) ClassifyScore() {
	if self.Score <= 0.5 {
		self.Classification = -1
	} else if self.Score <= 0.9 {
		self.Classification = 0
	} else {
		self.Classification = 1
	}
}

// Marshal the ObjMetadata struct into a byte array
func (self *TDIDF) Marshal() []byte {
	var bin bytes.Buffer

	// Create an encoder and send a value.
	enc := gob.NewEncoder(&bin)
	err := enc.Encode(self)

	if err != nil {
		log.Fatal("encode:", err)
	}

	return bin.Bytes()
}

// Marshal the User struct into a byte array
func (self *TDIDF) Unmarshal(u []byte) *TDIDF {
	dec := gob.NewDecoder(bytes.NewBuffer(u))

	err := dec.Decode(&self)

	if err != nil {
		log.Fatal("decode:", err)
	}

	return self
}
