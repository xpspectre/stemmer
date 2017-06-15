package stemmer

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"strings"
)

func TestIsVowel(t *testing.T) {
	letters := []rune{'a', 'b', 'c', 'y'}
	isVowels := []bool{true, false, false, true}
	for i, letter := range letters {
		assert.Equal(t, isVowels[i], IsVowel(letter))
	}
}

// See http://snowball.tartarus.org/texts/r1r2.html
func TestGetR1R2(t *testing.T) {
	words := []string{"beautiful", "beauty", "beau", "animadversion", "sprinkled", "eucharist"}
	r1s := []string{"iful", "y", "", "imadversion", "kled", "harist"}
	r2s := []string{"ul", "", "", "adversion", "", "ist"}
	for i, word := range words {
		r1, r2 := GetR1R2(word)
		assert.Equal(t, r1s[i], r1)
		assert.Equal(t, r2s[i], r2)
	}
}

func TestStem2Letter(t *testing.T) {
	words := []string{"a", "ax", "zy"}
	stems := words
	for i, word := range words {
		assert.Equal(t, stems[i], Stem(word))
	}
}

// http://snowball.tartarus.org/texts/vowelmarking.html
// The last example doesn't show up in english (yy) but illustrates the algorithm
func TestSetConsonantY(t *testing.T) {
	words := []string{"yes", "stay", "dyed", "ydyed", "ayyyyy"}
	newWords := []string{"Yes", "staY", "dyed", "Ydyed", "aYyYyY"}
	for i, word := range words {
		assert.Equal(t, newWords[i], SetConsonantY(word))
	}
}

func TestFindLongestSuffix(t *testing.T) {
	suffixes := []string{"'", "'s", "'s'"}
	words := []string{"there'", "there's", "there's'"}
	longesetSuffixes := []string{"'", "'s", "'s'"}
	newWords := []string{"there", "there", "there"}
	for i, word := range words {
		suffix := FindLongestSuffix(word, suffixes)
		assert.Equal(t, longesetSuffixes[i], suffix)
		assert.Equal(t, newWords[i], strings.TrimSuffix(word, suffix))
	}
}