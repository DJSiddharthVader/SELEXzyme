package main

func main() {
    var s Sequence
    s.seq = "AGCTCTCGGATCGATCGATATTTTAAAAAAA"
    s.fitness = 0.0
    var t Sequence
    t.seq = "TTTTAAAAAAAAGCTCTCGGATCGATCGATA"
    t.fitness = 0.0
    s.Complementarity(t)
}
