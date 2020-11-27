package main

import(
    "github.com/biogo/biogo/alphabet"
    "github.com/biogo/biogo/align"
)

//Constants
var DNA_ALPHABET = [4]rune{'A','C','G','T'}
var DNA_COMPLEMENTS = map[rune]rune{'A':'T','C':'G','G':'C','T':'A'}
var ALPHABET = alphabet.DNAgapped //alphabet for sequences
var SW_MATRIX = align.SWAffine { //alignment matrix for SW, example from biogo docs
        Matrix:  [][]int{        //       -   A   C   G   T
            {0, -1, -1, -1, -1}, // -     0  -1  -1  -1  -1
            {-1, 1, -2, -2, -2}, // A    -1   1  -2  -2  -2
            {-1, -2, 1, -2, -2}, // C    -1  -2   1  -2  -2
            {-1, -2, -2, 1, -2}, // G    -1  -2  -2   1  -2
            {-1, -2, -2, -2, 1}, // T    -1  -2  -2  -2   1
        }, GapOpen: -5, //gap opening penalty
}


//Data Types
type Member struct {
    seq string
    fitness float64
    label int
}
type Population []Member
//for getting alignment score from biogo
type Scorer interface {
    Score() int
}
