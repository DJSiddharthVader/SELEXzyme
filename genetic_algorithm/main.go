package main

func main() {
    /* // globals
    var DNA_ALPHABET = []string{"A","G","C","T"}
    var COMPLEMENT_MAP = map[rune]rune{'A':'T','T':'A','G':'C','C':'G'} //complementarty base pairs
    MINIMUM_HAIRPIN_LENGTH_PROPORTION := 0.5
    */
    // examples
    var s Sequence
    s.seq = "AGCTCTCGGATCGATCGATATTTTAAAAAAA"
    s.fitness = 0.0
    var t Sequence
    t.seq = "TTTTAAAAAAAAGCTCTCGGATCGATCGATA"
    t.fitness = 0.0
    // testing
    s.Complementarity(t)
    s.ReverseComplement()
}
