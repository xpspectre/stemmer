// Devel and test framework
// Read in sample vocab + stems and try out the code

package main

import (
	"fmt"
	"xpspectre.org/stemmer"
	"io/ioutil"
	"log"
	"strings"
)

func main() {
	// Load vocab + stems
	vocFile := "voc.txt"
	voc := make(map[string]string)
	data, err := ioutil.ReadFile(vocFile)
	if err != nil {
		log.Fatalf("Couldn't load voc file: %v\n", err)
	}
	for _, line := range strings.Split(string(data), "\n") {
		words := strings.Fields(line)
		if len(words) != 2 {
			continue
		}
		voc[words[0]] = words[1]
	}
	
	// Test out stemmer
	//	Keep track of mistakes
	totalWords := len(voc)
	wrongWords := 0
	for word, stem := range voc {
		stem_i := stemmer.Stem(word)
		if stem_i != stem {
			fmt.Printf("%s -> %s doesn't match %s\n", word, stem_i, stem)
			wrongWords++
		}
	}
	correctWords := totalWords - wrongWords
	pctCorrect := float64(correctWords) / float64(totalWords)
	fmt.Printf("%d/%d words stemmed correctly = %f%%\n", correctWords, totalWords, pctCorrect)
}
