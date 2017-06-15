// Implementation of the Porter stemming algorithm
// Source: https://tartarus.org/martin/PorterStemmer/
// *Also see: http://snowball.tartarus.org/algorithms/english/stemmer.html
// This is for English, assuming words are fed into it so ASCII only (probably)
// Assume inputs words to Stem are lowercase

package stemmer

import (
	"strings"
	"sort"
)

const vowels = "aeiuoy"
var double = [...]string{"bb", "dd", "ff", "gg", "mm", "nn", "pp", "rr", "tt"}
const liend = "cdeghkmnrt"

// Sort strings by length interfaces
//	Modified from https://gobyexample.com/sorting-by-functions
type ByDecLength []string
func (s ByDecLength) Len() int {
	return len(s)
}
func (s ByDecLength) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByDecLength) Less(i, j int) bool {
	return len(s[i]) > len(s[j]) // reverse order for decreasing sort
}

func Stem(s string) (string)  {
	// Words of <= 2 letters - leave as is
	if len(s) <= 2 {
		return s
	}
	
	// Remove initial apostrophe
	s = strings.TrimPrefix(s, "'")
	
	// Set consonant y's - initial y or y after a vowel
	// 	Denote consonants as capitalized Y
	s = SetConsonantY(s)
	
	// Step 0
	suffix := FindLongestSuffix(s, []string{"'", "'s", "'s'"})
	s = strings.TrimSuffix(s, suffix)
	
	return s
}

func IsVowel(c rune) (bool) {
	return strings.ContainsAny(string(c), vowels)
}

func SetConsonantY(s string) (string) {
	r := []rune(s)
	prevIsVowel := false
	for i, c := range s {
		if i == 0 && c == 'y' {
			r[0] = 'Y'
			continue // prevIsVowel = false implicitly
		}
		if prevIsVowel && c == 'y' {
			r[i] = 'Y'
			prevIsVowel = false
			continue
		}
		if IsVowel(c) {
			prevIsVowel = true
		}
	}
	return string(r)
}

// Search among list of suffixes. Returns longest suffix or empty string "" if none are found.
func FindLongestSuffix(s string, suffixes []string) (string) {
	sort.Sort(ByDecLength(suffixes))
	for _, suffix := range suffixes {
		if strings.HasSuffix(s, suffix) {
			return suffix
		}
	}
	return ""
}

// http://snowball.tartarus.org/texts/r1r2.html
// R1 is the region after the 1st non-vowel following a vowel, or null region at the end of the word if there isn't	a non-vowel
// R2 is the region after the 1st non-vowel following a vowel in R1, or null region at the end of the word if there isn't a non-vowel
func GetR1R2(s string) (string, string) {
	R1 := GetR1R2End(s)
	R2 := GetR1R2End(R1)
	return R1, R2
}

// Getting R1 and R2 is just applying the same procedure
func GetR1R2End(s string) (string) {
	// Find initial vowels. Start as consonant; then find 1st vowel; then find 1st consonant after that.
	initialVowel := false
	R1start := len(s) - 1
	Label:
	for i, c := range s {
		switch initialVowel {
		case false: // until 1st vowel
			if IsVowel(c) {
				initialVowel = true
			}
			continue
		case true: // until 1st consonant after vowel
			if !IsVowel(c) {
				R1start = i
				break Label  // break defaults to breaking out of switch case
			}
		}
	}
	return s[R1start+1:]
}

