package main

import (
	"bufio"
	// "bytes"
	"encoding/gob"
	"flag"
	"fmt"
	"github.com/zerklabs/sherlock"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

var (
	doc string

	fileToSearch = flag.String("file", "", "File to read in and search")

	isMinSCORESet = false
	isMinTDIDFSet = false

	lines []string

	maxRawFrequency int

	maxSCORE float64
	maxTDIDF float64
	minSCORE float64
	minTDIDF float64

	outputFileName    string
	outputGobFileName string
	outputHandle      *os.File
	outputGobHandle   *os.File

	printGarbage    = flag.Bool("garbage", false, "If enabled, will print low ranking lines")
	printSupervised = flag.Bool("supervised", false, "If enabled, will print lines needing supervision")

	tdidf []sherlock.TDIDF

	totalWords int
)

func main() {
	flag.Parse()

	if len(*fileToSearch) == 0 {
		log.Fatal("File expected to be given. Use --help")
	}

	outputFileName = fmt.Sprintf("runs/%d.out", time.Now().Unix())
	outputGobFileName = fmt.Sprintf("%s.gob", outputFileName)

	outputHandle, err := os.Create(outputFileName)
	defer outputHandle.Close()

	rfile, err := os.Open(*fileToSearch)
	defer rfile.Close()

	outputGobHandle, err := os.Create(outputGobFileName)
	defer outputGobHandle.Close()

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
				f := &sherlock.TDIDF{Word: v}
				tdidf = append(tdidf, *f)
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
	log.Printf("Maximum TD-IDF: %f\n", maxTDIDF)
	log.Printf("Minimum TD-IDF: %f\n", minTDIDF)
	log.Printf("Maximum Score: %f\n", maxSCORE)
	log.Printf("Minimum Score: %f\n", minSCORE)
	log.Printf("Wrote output to: %s\n", outputFileName)

	writeCollection()
	log.Printf("Wrote collection to: %s.gob\n", outputFileName)
}

func writeCollection() {
	pr, pw := io.Pipe()
	defer pr.Close()
	defer pw.Close()

	io.Copy(outputGobHandle, pr)

	// Create an encoder and send a value.
	enc := gob.NewEncoder(pw)
	err := enc.Encode(tdidf)

	if err != nil {
		log.Fatal("encode:", err)
	}
}

func assignNeighbors() {
	log.Println("Assigning neighboors")

	// assign neighbors
	var prevOne *sherlock.TDIDF
	var nextOne *sherlock.TDIDF

	for k, v := range tdidf {
		if k > 0 {
			prevOne = &tdidf[k-1]

			if k+1 < len(tdidf) {
				nextOne = &tdidf[k+1]
			}
		}

		v.LeftNeighbor = prevOne
		v.RightNeighbor = nextOne
		tdidf[k] = v
	}
}

// generate initial scores
func generateScores() {
	log.Println("Generating k-NN scores")

	for k, v := range tdidf {
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

		tdidf[k] = v
	}
}

func generateRawTF() {
	log.Println("Calculating Raw TF")
	for k, v := range tdidf {
		totalWords += 1
		v.RawFrequency(&doc)

		if v.RawTermFrequency > maxRawFrequency {
			maxRawFrequency = v.RawTermFrequency
		}

		tdidf[k] = v
	}
}

func generateTFIDF() {
	log.Println("Calculating TF-IDF")
	for k, v := range tdidf {
		v.Frequency(maxRawFrequency, totalWords)
		outputHandle.WriteString(fmt.Sprintf("%s|TDIDF:%f,TF:%f,IDF:%f,S:%f,C:%d\n", v.Word, v.Total, v.TermFrequency, v.InverseDocumentFrequency, v.Score, v.Classification))

		if v.Total > maxTDIDF {
			maxTDIDF = v.Total
		}

		if isMinTDIDFSet {
			if v.Total < minTDIDF {
				minTDIDF = v.Total
			}
		} else {
			minTDIDF = v.Total
			isMinTDIDFSet = true
		}

		tdidf[k] = v
	}
}
