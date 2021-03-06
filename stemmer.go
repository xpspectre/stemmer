// Implementation of the Porter stemming algorithm
// Source: https://tartarus.org/martin/PorterStemmer/
// *Also see: http://snowball.tartarus.org/algorithms/english/stemmer.html
// This is for English, assuming words are fed into it so ASCII only (probably)
// Assume inputs words to Stem are lowercase
// A smarter stemmer would return metadata about what the removed suffixes did
//	Has this improved performance in any testing?

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
	"ation": "ate",
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

// Do the following subs immediately and return
var exceptions = map[string]string{
	"skis": "ski", // start special changes
	"skies": "sky",
	"dying": "die",
	"lying": "lie",
	"tying": "tie",
	"idly": "idl", // start special -ly cases
	"gently": "gentl",
	"ugly": "ugli",
	"early": "earli",
	"only": "onli",
	"singly": "singl",
	"sky": "sky", // start invariant forms
	"news": "news",
	"howe": "howe",
	"atlas": "atlas", // looks like plural but not
	"cosmos": "cosmos",
	"bias": "bias",
	"andes": "andes",
}

// Leave these invariant after Step1a and return
var exceptions2 = map[string]string{
	"inning": "inning",
	"outing": "outing",
	"canning": "canning",
	"herring": "herring",
	"earring": "earring",
	"proceed": "proceed",
	"exceed": "exceed",
	"succeed": "succeed"}

// Sort strings by decreasing length interfaces
//	Modified from https://gobyexample.com/sorting-by-functions
type byDecLength []string
func (s byDecLength) Len() int {
	return len(s)
}
func (s byDecLength) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byDecLength) Less(i, j int) bool {
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
	
	// Exceptions directly mapped
	if stemmed, ok := exceptions[s]; ok {
		return stemmed
	}
	
	// Remove initial apostrophe
	s = strings.TrimPrefix(s, "'")
	
	s = setConsonantY(s)
	
	s = step0(s)
	s = step1a(s)
	
	// 2nd set of exceptions after trimming last 's'
	if stemmed, ok := exceptions2[s]; ok {
		return stemmed
	}
	
	s = step1b(s)
	s = step1c(s)
	s = step2(s)
	s = step3(s)
	s = step4(s)
	s = step5(s)
	
	// Convert consonanted Y back to y
	s = strings.Replace(s , "Y", "y", -1)
	
	return s
}

func step0(s string) (string) {
	suffix := findLongestSuffix(s, []string{"'", "'s", "'s'"})
	return strings.TrimSuffix(s, suffix)
}

func step1a(s string) (string) {
	//	The last "us" and "ss" are kept because they tell you to do nothing even though it ends in "s"
	suffix := findLongestSuffix(s, []string{"sses", "ied", "ies", "s", "us", "ss"})
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
		if len(s_) == 0 {
			return s
		}
		for _, c := range(s_[:len(s_)-1]) { // check all before 1 before suffix
			if isVowel(c) {
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

func step1b(s string) (string)  {
	suffix := findLongestSuffix(s, []string{"eed", "eedly", "ed", "edly", "ing", "ingly"})
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
			if isVowel(c) {
				suffix_ := findLongestSuffix(s_, []string{"at", "bl", "iz"})
				if suffix_ != "" {
					return s_ + "e"
				}
				suffix_ = findLongestSuffix(s_, double)
				if suffix_ != "" {
					r := []rune(s_)
					return string(r[:len(s_)-1])
				}
				if isShortWord(s_) {
					return s_ + "e"
				}
				return s_ // fell thru, no additional stuff after deleting suffix
			}
		}
		return s // fell thru, no vowels in preceding word part
	default: // do nothing
		return s
	}
	return s
}

func step1c(s string) (string) {
	r := []rune(s)
	rLen := len(r)
	if rLen < 2 {
		return s
	}
	rLast := r[rLen-1]
	rNext := r[rLen-2]
	// is the last check below redundant w/ the 2-letter ignore at the very top?
	if (rLast == 'y' || rLast == 'Y') && !isVowel(rNext) && rLen > 2 {
		r[rLen-1] = 'i'
	}
	return string(r)
}

func step2(s string) (string){
	R1 := GetR1(s)
	suffix := findLongestSuffix(s, step2Slice)
	if suffix != "" && strings.HasSuffix(R1, suffix) {
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
			if isLiEnding(r[len(r)-1]) {
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

func step3(s string) (string) {
	R1 := GetR1(s)
	suffix := findLongestSuffix(s, step3Slice)
	if suffix != "" && strings.HasSuffix(R1, suffix) {
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
			return s_ + step3Map[suffix]
		}
	}
	return s
}

// See https://tartarus.org/martin/PorterStemmer/, Common errors
//	The suffix is searched for in the original string, the longest is taken (if any are found),
//	and then the test (in R2) is applied. The in R2 test in applied only once.
func step4(s string) (string) {
	_, R2 := GetR1R2(s)
	suffix := findLongestSuffix(s, step4Words)
	if suffix != "" && strings.HasSuffix(R2, suffix) {
		return strings.TrimSuffix(s, suffix)
	}
	if strings.HasSuffix(R2, "ion") {
		s_ := strings.TrimSuffix(s, "ion")
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

func step5(s string) (string) {
	if strings.HasSuffix(s, "e") {
		R1, R2 := GetR1R2(s)
		if strings.HasSuffix(R2, "e") {
			return strings.TrimSuffix(s, "e")
		}
		if strings.HasSuffix(R1, "e") {
			s_ := strings.TrimSuffix(s, "e")
			if !isEndShortSyllable(s_) {
				return s_
			} else{
				return s
			}
		}
		return s
	}
	if strings.HasSuffix(s, "l") {
		_, R2 := GetR1R2(s)
		s_ := strings.TrimSuffix(s, "l")
		r := []rune(s_)
		rLast := r[len(r)-1]
		if strings.HasSuffix(R2, "l") && rLast == 'l' {
			return s_
		} else {
			return s
		}
	}
	return s
}

func isVowel(c rune) (bool) {
	return strings.ContainsAny(string(c), vowels)
}

func isLiEnding(c rune) (bool) {
	return strings.ContainsAny(string(c), liend)
}

// Search ends of words only
// Check only/start syllable for short words; else the 3-letter suffix for long words
func isEndShortSyllable(s string) (bool) {
	r := []rune(s)
	//
	if len(r) < 3 {
		if isVowel(r[0]) && !isVowel(r[1]) {
			return true
		} else {
			return false
		}
	} else {
		r = r[len(r)-3:]
		if !isVowel(r[0]) && isVowel(r[1]) && !isVowel(r[2]) && r[2] != 'w' && r[2] !='x' && r[2] != 'Y' {
			return true
		} else {
			return false
		}
	}
}

// Short word ends in a short syllable and R1 is null
func isShortWord(s string)(bool) {
	R1 := GetR1(s)
	if isEndShortSyllable(s) && R1 == "" {
		return true
	}
	return false
}

// Set consonant y's - initial y or y after a vowel
// 	Denote consonants as capitalized Y
func setConsonantY(s string) (string) {
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
		if isVowel(c) {
			prevIsVowel = true
		} else {
			prevIsVowel = false
		}
	}
	return string(r)
}

// Search among list of suffixes. Returns longest suffix or empty string "" if none are found.
func findLongestSuffix(s string, suffixes []string) (string) {
	sort.Sort(byDecLength(suffixes))
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

// Overstemming can result from these prefixes
var r1Exceptions = []string{"gener", "commun", "arsen"}

func GetR1(s string) (string) {
	// Handle 3 corner cases
	for i := 0; i < len(r1Exceptions); i++ {
		if strings.HasPrefix(s, r1Exceptions[i]) {
			return strings.TrimPrefix(s, r1Exceptions[i])
		}
	}
	return getR1R2End(s)
}

func GetR1R2(s string) (string, string) {
	R1 := GetR1(s)
	return R1, getR1R2End(R1)
}

// Getting R1 and R2 is just applying the same procedure
func getR1R2End(s string) (string) {
	// Find initial vowels. Start as consonant; then find 1st vowel; then find 1st consonant after that.
	initialVowel := false
	R1start := len(s) - 1 // default for specifying null region at the end
	Label:
	for i, c := range s {
		switch initialVowel {
		case false: // until 1st vowel
			if isVowel(c) {
				initialVowel = true
			}
			continue
		case true: // until 1st consonant after vowel
			if !isVowel(c) {
				R1start = i
				break Label  // break defaults to breaking out of switch case
			}
		}
	}
	return s[R1start+1:]
}

