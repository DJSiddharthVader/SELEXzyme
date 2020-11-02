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
// input: a population of sequences and how many you will pick (proportion is in (0,1)
// output: new population of Sequences
func BreedNewGeneration(generation Population) Population {
    nextGeneration := make(Population,len(generation))
    fittestMembers := GetFittestMembers(generation)
    for i,member := range fittestMembers {
        member.label = i
        nextGeneration[i] = Member{label:i,
                                     fitness:member.fitness,
                                     seq:member.seq,
                                    }
    }
    for i:=len(fittestMembers);i<len(nextGeneration);i++ {
        //breed new sequences untill our new generation is same size as previous
        nextGeneration[i] = BreedSequence(fittestMembers,i)
    }
    return nextGeneration
}

//Decide when to stop breeding new generations
// FitnessPlateau() checks if a population has plateaued in fitness
// actually checks if covarinace of mean fitness of the last n generations is < tolerance
// n and toleracen are user defined
// input: population
// output: whether fitness has plateued, bool
func FitnessPlateau(generationFitnesses []float64) bool {
// const FITNESS_PLATEAU_GENERATIONS = 5 //number of generations for which fintess must have plateaued to halt simulation
    cov := CoV(generationFitnesses[FITNESS_PLATEAU_GENERATIONS:len(generationFitnesses)])
    return (cov < FITNESS_PLATEAU_TOLERANCE)
}
// RunSimulation() runs a genetic algorithm and returns the final generation
// input: number of seqs in poplation, upper and lower bound sequence length
// for the initial population, maximum number of iterations
// output: last generation of sequences
func RunSimulation(size,lower,upper,maxIterations int) Population {
    generationFitnesses := make([]float64,maxIterations)
    currentGen := InitializeGeneration(size,lower,upper)
    for gen := 0; gen < maxIterations; gen++ {//terminate regardless after maxIterations
        generationFitnesses = append(generationFitnesses,currentGen.MeanFitness())
        if FitnessPlateau(generationFitnesses) {//check if average fitness has plateaued
            return currentGen //if plateau, no improvements from continnuing simulation, finish
        }
        currentGen = BreedNewGeneration(currentGen)
    }//never reached fitness plateau
    return currentGen
}

