package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/zerklabs/sherlock"
	"log"
	"os"
	"strings"
)

var (
	doc string

	fileToSearch = flag.String("file", "", "File to read in and search")

	isMinSCORESet = false
	isMinTFIDFSet = false

	lines []string

	maxRawFrequency int

	maxSCORE float64
	maxTFIDF float64
	minSCORE float64
	minTFIDF float64

	printGarbage    = flag.Bool("garbage", false, "If enabled, will print low ranking lines")
	printSupervised = flag.Bool("supervised", false, "If enabled, will print lines needing supervision")

	tfidf []sherlock.TFIDF

	totalWords int
)

func main() {
	flag.Parse()

	if len(*fileToSearch) == 0 {
		log.Fatal("File expected to be given. Use --help")
	}

	rfile, err := os.Open(*fileToSearch)
	defer rfile.Close()

	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(rfile)
	for scanner.Scan() {
		text := scanner.Text()

		lines = append(lines, text)
		doc = strings.Join(lines, "\n")
		words := sherlock.TokenizeLine(text)

		// create initial collection
		for _, v := range words {
			if len(v) > 0 {
				f := &sherlock.TFIDF{Word: v}
				tfidf = append(tfidf, *f)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}

	log.Println("File completely read in")

	assignNeighbors()
	generateScores()
	generateRawTF()
	generateTFIDF()

	log.Printf("Total words: %d\n", totalWords)
	log.Printf("Total lines: %d\n", len(lines))
	log.Printf("Maximum TD-IDF: %f\n", maxTFIDF)
	log.Printf("Minimum TD-IDF: %f\n", minTFIDF)
	log.Printf("Maximum Score: %f\n", maxSCORE)
	log.Printf("Minimum Score: %f\n", minSCORE)

	// writeCollection()
	// log.Printf("Wrote collection to: %s.gob\n", outputFileName)
}

// func writeCollection() {
// 	pr, pw := io.Pipe()
// 	defer pr.Close()
// 	defer pw.Close()

// 	io.Copy(outputGobHandle, pr)

// 	// Create an encoder and send a value.
// 	enc := gob.NewEncoder(pw)
// 	err := enc.Encode(tfidf)

// 	if err != nil {
// 		log.Fatal("encode:", err)
// 	}
// }

func assignNeighbors() {
	log.Println("Assigning neighboors")

	// assign neighbors
	var prevOne *sherlock.TFIDF
	var nextOne *sherlock.TFIDF

	for k, v := range tfidf {
		if k > 0 {
			prevOne = &tfidf[k-1]

			if k+1 < len(tfidf) {
				nextOne = &tfidf[k+1]
			}
		}

		v.LeftNeighbor = prevOne
		v.RightNeighbor = nextOne
		tfidf[k] = v
	}
}

// generate initial scores
func generateScores() {
	log.Println("Generating k-NN scores")

	for k, v := range tfidf {
		v.ScoreWord()
		v.ClassifyScore()

		if v.Score > maxSCORE {
			maxSCORE = v.Score
		}

		if isMinSCORESet {
			if v.Score < minSCORE {
				minSCORE = v.Score
			}
		} else {
			minSCORE = v.Score
			isMinSCORESet = true
		}

		tfidf[k] = v
	}
}

func generateRawTF() {
	log.Println("Calculating Raw TF")

	for k, v := range tfidf {
		totalWords += 1
		v.RawFrequency(&doc)

		if v.RawTermFrequency > maxRawFrequency {
			maxRawFrequency = v.RawTermFrequency
		}

		tfidf[k] = v
	}
}

func generateTFIDF() {
	log.Println("Calculating TF-IDF")
	for k, v := range tfidf {
		v.Frequency(maxRawFrequency, totalWords)
		tfidf[k] = v

		fmt.Printf("%s|TFIDF:%f,TF:%f,IDF:%f,S:%f,C:%d\n", v.Word, v.Total, v.TermFrequency, v.InverseDocumentFrequency, v.Score, v.Classification)

		if v.Total > maxTFIDF {
			maxTFIDF = v.Total
		}

		if isMinTFIDFSet {
			if v.Total < minTFIDF {
				minTFIDF = v.Total
			}
		} else {
			minTFIDF = v.Total
			isMinTFIDFSet = true
		}

	}
}
