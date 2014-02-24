package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/zerklabs/sherlock"
	"log"
	"os"
)

var printGarbage = flag.Bool("garbage", false, "If enabled, will print low ranking lines")
var printSupervised = flag.Bool("supervised", false, "If enabled, will print lines needing supervision")

func main() {
	flag.Parse()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		classify(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}

}

func classify(line string) {
	words := sherlock.TokenizeLine(line)

	for _, v := range words {
		score := sherlock.ScoreWord(v)
		classify := sherlock.ClassifyScore(score)

		if *printGarbage && classify == -1 {
			fmt.Printf("%d,%v,%s\n", classify, score, v)
		} else if *printSupervised && classify == 0 {
			fmt.Printf("%d,%v,%s\n", classify, score, v)
		} else if classify == 1 {
			fmt.Printf("%d,%v,%s\n", classify, score, v)
		}
	}
}
