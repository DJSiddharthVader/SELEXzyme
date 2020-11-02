package main

import(
    "fmt"
    "math/rand"
)
func newlinebreed() {
    fmt.Println()
}

//Initialize random pool of sequences
// MakeRandomSeq() returns a random DNA string of the given length
// input: int length of the sequence
// output: string DNA sequence
func MakeRandomSeq(length int) string {
    var seq string
    for i := 0; i < length; i++ {
        seq = seq + PickRandomBase() // pick a random element from the alphabet
    }
    return seq
}
// MakeRandomSequence() returns a random Sequence Object with a seq of the given length
// input: length of the sequence
// output: Sequence object
func MakeRandomSequence(length int) Member {
    var s Member
    s.seq = MakeRandomSeq(length)
    s.ScoreFitness()
    return s
}
// InitializeGeneration() create a random pool of sequences to start our gentic algorithmj
// input:  the number of sequences to generate and lower,upper bounds onsequence length
// output: a new random population (slice of Sequences) with size members
func InitializeGeneration(size,lower,upper int) Population {
    population := make([]Member,size)
    for i := 0; i < size; i++ {
        population[i] = MakeRandomSequence(RandomIntBetween(lower,upper))
    }
    return population
}

type Population []Sequence
/*
// SelectMembersForBreeding() select the best members from the current population for breeding the next generation
// BreedNewGeneration() create a new population from previous best members and breeding new members from them
// BreedUntilFinished() keep breeding new generations until average fitness plateaus or a specific number of generations has passed
// WriteToFile() write every member of the population into a file, either tsv or fasta
*/
