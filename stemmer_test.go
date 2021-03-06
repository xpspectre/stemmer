package stemmer

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestIsVowel(t *testing.T) {
	letters := []rune{'a', 'b', 'c', 'y'}
	isVowels := []bool{true, false, false, true}
	for i, letter := range letters {
		assert.Equal(t, isVowels[i], isVowel(letter))
	}
}

func TestIsLiEnding(t *testing.T) {
	letters := []rune{'a', 'b', 'c', 't'}
	isLiEnds := []bool{false, false, true, true}
	for i, letter := range letters {
		assert.Equal(t, isLiEnds[i], isLiEnding(letter))
	}
}

// See http://snowball.tartarus.org/texts/r1r2.html
func TestGetR1R2(t *testing.T) {
	words := []string{"beautiful", "beauty", "beau", "animadversion", "sprinkled", "eucharist"}
	r1s := []string{"iful", "y", "", "imadversion", "kled", "harist"}
	r2s := []string{"ul", "", "", "adversion", "", "ist"}
	for i, word := range words {
		r1 := GetR1(word)
		_, r2 := GetR1R2(word)
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
	words := []string{"yes", "stay", "dyed", "ydyed", "ayyyyy", "willingly"}
	newWords := []string{"Yes", "staY", "dyed", "Ydyed", "aYyYyY", "willingly"}
	for i, word := range words {
		assert.Equal(t, newWords[i], setConsonantY(word))
	}
}

func TestFindLongestSuffix(t *testing.T) {
	suffixes := []string{"'", "'s", "'s'"} // suffixes are substrings so the test must find longest
	words := []string{"there'", "there's", "there's'"}
	longesetSuffixes := []string{"'", "'s", "'s'"}
	for i, word := range words {
		assert.Equal(t, longesetSuffixes[i], findLongestSuffix(word, suffixes))
	}
}

func TestStep0(t *testing.T) {
	words := []string{"there'", "there's", "there's'"}
	newWords := []string{"there", "there", "there"}
	for i, word := range words {
		assert.Equal(t, newWords[i], step0(word))
	}
}

func TestStep1a(t *testing.T) {
	assert.Equal(t, "blah", step1a("blah"))
	assert.Equal(t, "fdass", step1a("fdasses"))
	assert.Equal(t, "tie", step1a("ties"))
	assert.Equal(t, "cri", step1a("cries"))
	assert.Equal(t, "gas", step1a("gas"))
	assert.Equal(t, "this", step1a("this"))
	assert.Equal(t, "gap", step1a("gaps"))
	assert.Equal(t, "kiwi", step1a("kiwis"))
}

func TestStep1b(t *testing.T) {
	assert.Equal(t, "airspee", step1b("airspeed")) // is this actually a good example?
	assert.Equal(t, "creed", step1b("creed"))
	assert.Equal(t, "luxuriate", step1b("luxuriatedly"))
	assert.Equal(t, "hop", step1b("hopping"))
	assert.Equal(t, "hope", step1b("hoping"))
	assert.Equal(t, "resol", step1b("resoled"))
}

func TestIsEndShortSyllable(t *testing.T) {
	endShort := []string{"rap", "trap", "entrap", "ow", "on", "at"}
	notEndShort := []string{"uproot", "bestow", "disturb"}
	for _, word := range endShort {
		assert.True(t, isEndShortSyllable(word))
	}
	for _, word := range notEndShort {
		assert.False(t, isEndShortSyllable(word))
	}
}

func TestIsShortWord(t *testing.T) {
	shortWords := []string{"bed", "shed", "shred"}
	notShortWords := []string{"bead", "embed", "beds"}
	for _, word := range shortWords {
		assert.True(t, isShortWord(word))
	}
	for _, word := range notShortWords {
		assert.False(t, isShortWord(word))
	}
}

func TestStep1c(t *testing.T) {
	words := []string{"cry", "by", "say"}
	newWords := []string{"cri", "by", "say"}
	for i, word := range words {
		assert.Equal(t, newWords[i], step1c(word))
	}
}

func TestStep2(t *testing.T) {
	words := []string{"additional", "relational", "yogi", "stimuli"}
	newWords := []string{"addition", "relate", "yogi", "stimuli"}
	for i, word := range words {
		assert.Equal(t, newWords[i], step2(word))
	}
}

func TestStep3(t *testing.T) {
	// couldn't think of a real word for the last one
	words := []string{"conditional", "procrastinative"}
	newWords := []string{"condition", "procrastin"}
	for i, word := range words {
		assert.Equal(t, newWords[i], step3(word))
	}
}

func TestStep4(t *testing.T) {
	words := []string{"vocalize", "materialize", "vocation", "petition", "basement"}
	newWords := []string{"vocal", "material", "vocat", "petit", "basement"}
	for i, word := range words {
		assert.Equal(t, newWords[i], step4(word))
	}
}

func TestStep5(t *testing.T) {

}

func TestStemExceptions(t *testing.T) {
	words := []string{"skis", "gently", "bias"}
	newWords := []string{"ski", "gentl", "bias"}
	for i, word := range words {
		assert.Equal(t, newWords[i], Stem(word))
	}
}

func TestStemExceptions2(t *testing.T) {
	words := []string{"earrings"}
	newWords := []string{"earring"}
	for i, word := range words {
		assert.Equal(t, newWords[i], Stem(word))
	}
}

// Overstemming corner case
//	Also useful for testing a bunch of the steps - making sure they refer to the right maps and such
func TestOverstem(t *testing.T) {
	words := []string{"generate", "generates", "generated", "generating",
		"general", "generally",
		"generic", "generically",
		"generous", "generously"}
	newWords := []string{"generat", "generat", "generat", "generat",
		"general", "general",
		"generic", "generic",
		"generous", "generous"}
	for i, word := range words {
		assert.Equal(t, newWords[i], Stem(word))
	}
}

// Basic tests to make sure the whole thing runs
func TestStem(t *testing.T) {
	words := []string{"the"}
	newWords := []string{"the"}
	for i, word := range words {
		assert.Equal(t, newWords[i], Stem(word))
	}
}