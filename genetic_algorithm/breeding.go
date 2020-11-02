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

//Breed new generation of Sequences
// Crossover() creates a new sequence by crossing over with some input sequence at a random locus and rescores the new sequence
// input: sequences to crossover
// output: new sequence that is a hybrid of the inputs
func (s Member) Crossover(t Member) string {
    crossOverIndex := RandomIntBetween(0,Min(len(s.seq),len(t.seq))) //where to crossover
    // combine front half of s.seq and back half of t.seq
    return s.seq[0:crossOverIndex] + t.seq[crossOverIndex:len(t.seq)]
}
// Mutate() mutates a DNA sequence at each position with some probability
// input: probability that each site will be mutated
// output: sequence with mutations
func (s Member) Mutate() string {
    var mutated string
    for _,base := range s.seq {
        if rand.Float64() > MUTATION_RATE {//dont mutate
            mutated = mutated + string(base)
        } else {//mutate this base to a new base
            mutated = mutated + PickDifferentRandomBase(base)
        }
    }
    return mutated
}
// GetFittestMembers() selects the fittest members from the current population
// for breeding the next generation
// input: a population of sequences and how many you will pick (proportion is in (0,1)
// output: the top proportion percent of the generation sequences by fitness
func GetFittestMembers(generation Population) Population {
    // the index at proportion percent of generation from the end
    index := len(generation) - int(float64(len(generation))*TOP_SEQUENCE_PERCENT)
    return generation.SortByFitness()[index:len(generation)]
}
// BreedSequence() breeds a new sequence from a population
// input: some set of sequences
// output: a single new sequence bred from 2 random population members
func BreedSequence(pop Population, label int) Member {
    seq1 := pop[rand.Intn(len(pop))] //pick a random Sequence
    seq2 := pop[rand.Intn(len(pop))] //pick another random Sequence
    newSequence := Member{seq:seq1.seq}
    newSequence.seq = newSequence.Crossover(seq2)
    newSequence.seq = newSequence.Mutate()
    newSequence.fitness = newSequence.ScoreFitness()
    newSequence.label = label
    return newSequence
}
// BreedNewGeneration() create a new population from previous best members and breeding new members from them
// BreedUntilFinished() keep breeding new generations until average fitness plateaus or a specific number of generations has passed
// WriteToFile() write every member of the population into a file, either tsv or fasta
*/
