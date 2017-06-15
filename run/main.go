// Devel and test framework

package main

import (
	"fmt"
	"xpspectre.org/stemmer"
)

func main()  {
	R1, R2 := stemmer.GetR1R2("generically")
	fmt.Println(R1)
	fmt.Println(R2)
	fmt.Println(stemmer.Stem("generically"))
	//fmt.Println(strings.HasSuffix("", "ing"))
}
