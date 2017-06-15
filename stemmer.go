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
var double = []string{"bb", "dd", "ff", "gg", "mm", "nn", "pp", "rr", "tt"}
const liend = "cdeghkmnrt"

var step2Map = map[string]string{
	"tional": "tion",
	"enci": "ence",
	"anci": "ance",
	"abli": "able",
	"entli": "ent",
	"izer": "ize",
	"ization": "ize",
	"ational": "ate",
	"ation": "age",
	"ator": "ate",
	"alism": "al",
	"aliti": "al",
	"alli" : "al",
	"fulness": "ful",
	"ousli": "ous",
	"ousness": "ous",
	"iveness": "ive",
	"iviti": "ive",
	"biliti": "ble",
	"bli": "ble",
	"fulli": "ful",
	"lessli": "less",
	"ogi": "og", // special handling
	"li": "", // special handling
}
var step2Slice = make([]string, len(step2Map)) // holds just the keys of above

var step3Map = map[string]string{
	"tional": "tion",
	"ational": "ate",
	"alize": "al",
	"icate": "ic",
	"iciti": "ic",
	"ical": "ic",
	"ful": "",
	"ness": "",
	"ative": "",
}
var step3Slice = make([]string, len(step3Map))

var step4Words = []string{"al", "ance", "ence", "er", "ic", "able", "ible", "ant", "ement", "ment", "ent", "ism",
	"ate", "iti", "ous", "ive", "ize"}

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

func init() {
	// Get slices of map keys for suffix lookups
	i := 0
	for suffix := range step2Map {
		step2Slice[i] = suffix
		i++
	}
	
	i = 0
	for suffix := range step3Map {
		step3Slice[i] = suffix
		i++
	}
}

func Stem(s string) (string)  {
	// Words of <= 2 letters - leave as is
	if len(s) <= 2 {
		return s
	}
	
	// Remove initial apostrophe
	s = strings.TrimPrefix(s, "'")
	
	s = SetConsonantY(s)
	
	s = Step0(s)
	s = Step1a(s)
	s = Step1b(s)
	s = Step1c(s)
	s = Step2(s)
	s = Step3(s)
	s = Step4(s)
	s = Step5(s)
	
	return s
}

func Step0(s string) (string) {
	suffix := FindLongestSuffix(s, []string{"'", "'s", "'s'"})
	return strings.TrimSuffix(s, suffix)
}

func Step1a(s string) (string) {
	//	The last "us" and "ss" are kept because they tell you to do nothing even though it ends in "s"
	suffix := FindLongestSuffix(s, []string{"sses", "ied", "ies", "s", "us", "ss"})
	switch {
	case suffix == "sses":
		return strings.TrimSuffix(s, "sses") + "ss"
	case suffix == "ied" || suffix == "ies":
		s_ := strings.TrimSuffix(s, suffix)
		if len(s_) > 1 {
			return s_ + "i"
		} else {
			return s_ + "ie"
		}
	case suffix == "s":
		s_ := strings.TrimSuffix(s, suffix)
		for _, c := range(s_[:len(s_)-1]) { // check all before 1 before suffix
			if IsVowel(c) {
				return s_
			}
		}
		return s // fell thru, all before 1 before suffix are consonants
	case suffix == "us" || suffix == "ss":
		return s
	default: // do nothing
		return s
	}
}

func Step1b(s string) (string)  {
	suffix := FindLongestSuffix(s, []string{"eed", "eedly", "ed", "edly", "ing", "ingly"})
	switch {
	case suffix == "eed" || suffix == "eedly":
		R1 := GetR1(s)
		if strings.HasSuffix(R1, suffix) {
			return strings.TrimSuffix(s, suffix) + "ee"
		} else {
			return s
		}
	case suffix == "ed" || suffix == "edly" || suffix == "ing" || suffix == "ingly":
		s_ := strings.TrimSuffix(s, suffix)
		for _, c := range s_ {
			if IsVowel(c) {
				suffix_ := FindLongestSuffix(s_, []string{"at", "bl", "iz"})
				if suffix_ != "" {
					return s_ + "e"
				}
				suffix_ = FindLongestSuffix(s_, double)
				if suffix_ != "" {
					r := []rune(s_)
					return string(r[:len(s_)-1])
				}
				if IsShortWord(s_) {
					return s_ + "e"
				}
			}
		}
		return s // fell thru, no vowels in preceding word part
	default: // do nothing
		return s
	}
	return s
}

func Step1c(s string) (string) {
	r := []rune(s)
	rLen := len(r)
	rLast := r[rLen-1]
	rNext := r[rLen-2]
	// is the last check below redundant w/ the 2-letter ignore at the very top?
	if (rLast == 'y' || rLast == 'Y') && !IsVowel(rNext) && rLen > 2 {
		r[rLen-1] = 'i'
	}
	return string(r)
}



func Step2(s string) (string){
	R1 := GetR1(s)
	suffix := FindLongestSuffix(R1, step2Slice)
	if suffix != "" {
		s_ := strings.TrimSuffix(s, suffix)
		switch suffix {
		case "ogi":
			r := []rune(s_)
			if r[len(r)-1] == 'l' {
				return s_ + "og"
			} else {
				return s
			}
		case "li":
			r := []rune(s_)
			if IsLiEnding(r[len(r)-1]) {
				return s_
			} else {
				return s
			}
		default:
			return s_ + step2Map[suffix]
		}
	}
	return s
}

func Step3(s string) (string) {
	R1 := GetR1(s)
	suffix := FindLongestSuffix(R1, step3Slice)
	if suffix != "" {
		s_ := strings.TrimSuffix(s, suffix)
		switch suffix {
		case "ative":
			_, R2 := GetR1R2(s)
			if strings.HasSuffix(R2, "ative") {
				return s_
			} else {
				return s
			}
		default:
			return s_ + step2Map[suffix]
		}
	}
	return s
}

func Step4(s string) (string) {
	_, R2 := GetR1R2(s)
	suffix := FindLongestSuffix(R2, step4Words)
	if suffix != "" {
		return strings.TrimSuffix(s, suffix)
	}
	if strings.HasSuffix(R2, "ion") {
		s_ := strings.TrimSuffix(s, suffix)
		r := []rune(s_)
		rLast := r[len(r)-1]
		if rLast == 's' || rLast == 't' {
			return s_
		} else {
			return s
		}
	}
	return s
}

func Step5(s string) (string) {
	return s
}

func IsVowel(c rune) (bool) {
	return strings.ContainsAny(string(c), vowels)
}

func IsLiEnding(c rune) (bool) {
	return strings.ContainsAny(string(c), liend)
}

// Search ends of words only
// Check only/start syllable for short words; else the 3-letter suffix for long words
func IsEndShortSyllable(s string) (bool) {
	r := []rune(s)
	//
	if len(r) < 3 {
		if IsVowel(r[0]) && !IsVowel(r[1]) {
			return true
		} else {
			return false
		}
	} else {
		r = r[len(r)-3:]
		if !IsVowel(r[0]) && IsVowel(r[1]) && !IsVowel(r[2]) && r[2] != 'w' && r[2] !='x' && r[2] != 'Y' {
			return true
		} else {
			return false
		}
	}
}

// Short word ends in a short syllable and R1 is null
func IsShortWord(s string)(bool) {
	R1 := GetR1(s)
	if IsEndShortSyllable(s) && R1 == "" {
		return true
	}
	return false
}

// Set consonant y's - initial y or y after a vowel
// 	Denote consonants as capitalized Y
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
func GetR1(s string) (string) {
	return GetR1R2End(s)
}

func GetR1R2(s string) (string, string) {
	R1 := GetR1(s)
	return R1, GetR1R2End(R1)
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

