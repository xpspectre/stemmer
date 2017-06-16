# stemmer

English Porter2 stemming algorithm  in Go. Mostly for learning Go, strings, and stemming.

Exposes the function `func Stem(s string) (string)` which gives you the stem for some word.

Run unit tests in `stemmer_test.go`. Run `main/main.go` to test the stemmer on a big vocab with precomputed stems. Used `dev/main.go` for misc debugging of particular words.