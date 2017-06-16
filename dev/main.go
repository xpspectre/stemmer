package main

import (
	"xpspectre.org/stemmer"
	"fmt"
)

func main()  {
	word := "basement"
	stem := "base"
	stem_ := stemmer.Stem(word)
	
	R1, R2 := stemmer.GetR1R2(word)
	
	fmt.Printf("%s / %s\n", R1, R2)
	fmt.Printf("%s -> %s / %s\n", word, stem_, stem)
}
