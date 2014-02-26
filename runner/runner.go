package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/zerklabs/sherlock"
	"log"
	"os"
	"strings"
	"time"
)

var (
	printGarbage    = flag.Bool("garbage", false, "If enabled, will print low ranking lines")
	printSupervised = flag.Bool("supervised", false, "If enabled, will print lines needing supervision")
	fileToSearch    = flag.String("file", "", "File to read in and search")
	lines           []string
	doc             string
	totalWords      int
	tdidf           []sherlock.TDIDF
	storageChan     chan *sherlock.TDIDF
	outputHandle    *os.File

	rawDone         int
	totDone         int
	maxRawFrequency int

	maxTDIDF      float64
	minTDIDF      float64
	maxSCORE      float64
	minSCORE      float64
	isMinTDIDFSet = false
	isMinSCORESet = false
)

func main() {
	flag.Parse()

	if len(*fileToSearch) == 0 {
		log.Fatal("File expected to be given. Use --help")
	}

	outputHandle, err := os.Create(fmt.Sprintf("runs/%d-tdidf.out", time.Now().Unix()))
	defer outputHandle.Close()

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

		for _, v := range words {
			if len(v) > 0 {
				f := &sherlock.TDIDF{Word: v}
				f.ScoreWord()
				f.ClassifyScore()

				if f.Score > maxSCORE {
					maxSCORE = f.Score
				}

				if isMinSCORESet {
					if f.Score < minSCORE {
						minSCORE = f.Score
					}
				} else {
					minSCORE = f.Score
					isMinSCORESet = true
				}

				tdidf = append(tdidf, *f)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}

	log.Println("File completely read in")
	log.Println("Calculating Raw TF")
	for k, v := range tdidf {
		totalWords += 1
		v.RawFrequency(&doc)

		tdidf[k] = v

		if v.RawTermFrequency > maxRawFrequency {
			maxRawFrequency = v.RawTermFrequency
		}
	}

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

	log.Printf("Total words: %d\n", totalWords)
	log.Printf("Total lines: %d\n", len(lines))
	log.Printf("Maximum TD-IDF: %f\n", maxTDIDF)
	log.Printf("Minimum TD-IDF: %f\n", minTDIDF)
	log.Printf("Maximum Score: %f\n", maxSCORE)
	log.Printf("Minimum Score: %f\n", minSCORE)
}
