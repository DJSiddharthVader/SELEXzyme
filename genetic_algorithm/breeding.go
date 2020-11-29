package main

import(
    "fmt"
    "math/rand"
    "github.com/biogo/biogo/seq/linear"
    "github.com/cheggaaa/pb"
)

const indel_rate = 0.1 // probability a mutation is an indel

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
func MakeRandomSequence(length int, target *linear.Seq) Member {
    var s Member
    s.seq = MakeRandomSeq(length)
    s.ScoreFitness(target)
    return s
}
// InitializeGeneration() create a random pool of sequences to start our gentic algorithm
// input:  the number of sequences to generate and lower,upper bounds onsequence length
// output: a new random population (slice of Sequences) with size members
func InitializeGeneration(size,lower,upper int,target *linear.Seq) Population {
    population := make([]Member,size)
    for i := 0; i < size; i++ {
        population[i] = MakeRandomSequence(RandomIntBetween(lower,upper),target)
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
func (s Member) Mutate(mutation_rate,indel_rate float64) string {
    var mutated string
    for _,base := range s.seq {
        if rand.Float64() > mutation_rate {//dont mutate
            mutated = mutated + string(base)
        } else {//mutate this base to a new base
            if rand.Float64() > indel_rate {//regular mutation
                mutated = mutated + PickDifferentRandomBase(base)
            } else {
                if rand.Intn(2) == 0 {// insert new base
                    mutated = mutated + string(base) + PickRandomBase()
                } else {
                    //delete base (add nothing)
                }
            }
        }
    }
    return mutated
}
// GetFittestMembers() selects the fittest members from the current population
// for breeding the next generation
// input: a population of sequences and how many you will pick (proportion is in (0,1)
// output: the top proportion percent of the generation sequences by fitness
func GetFittestMembers(generation Population, top_sequence_percent float64) Population {
    // the index at proportion percent of generation from the end
    index := len(generation) - int(float64(len(generation))*top_sequence_percent)
    return generation.SortByFitness()[index:len(generation)]
}
// BreedSequence() breeds a new sequence from a population
// input: some set of sequences
// output: a single new sequence bred from 2 random population members
func BreedSequence(pop Population, label int,target *linear.Seq, mutation_rate,indel_rate float64) Member {
    seq1 := pop[rand.Intn(len(pop))] //pick a random Sequence
    seq2 := pop[rand.Intn(len(pop))] //pick another random Sequence
    newSequence := Member{seq:seq1.seq}
    newSequence.seq = newSequence.Crossover(seq2)
    newSequence.seq = newSequence.Mutate(mutation_rate,indel_rate)
    newSequence.fitness = newSequence.ScoreFitness(target)
    newSequence.label = label
    return newSequence
}
// BreedNewGeneration() create a new population from previous best members and breeding new members from them
// input: a population of sequences and how many you will pick (proportion is in (0,1)
// output: new population of Sequences
func BreedNewGeneration(generation Population, target *linear.Seq, mutation_rate float64, indel_rate float64, top_sequence_percent float64) Population {
    nextGeneration := make(Population,len(generation))
    fittestMembers := GetFittestMembers(generation,top_sequence_percent)
    for i,member := range fittestMembers {
        member.label = i
        nextGeneration[i] = Member{label:i,
                                     fitness:member.fitness,
                                     seq:member.seq,
                                    }
    }
    for i:=len(fittestMembers);i<len(nextGeneration);i++ {
        //breed new sequences untill our new generation is same size as previous
        nextGeneration[i] = BreedSequence(fittestMembers,i,target,mutation_rate,indel_rate)
    }
    return nextGeneration
}

//Decide when to stop breeding new generations
// FitnessPlateau() checks if a population has plateaued in fitness
// actually checks if covarinace of mean fitness of the last n generations is < tolerance
// n and toleracen are user defined
// input: mean fitness of previous generations
// output: whether fitness has plateued, bool
func FitnessPlateau(mode string, fitnesses [][]float64, fitness_plateau_tolerance float64) bool {
    switch mode {
        case "cov_mean"://mean fitness CoV for last few generations
            var meanFitnesses []float64
            for _,fitness := range fitnesses {
                meanFitnesses = append(meanFitnesses,Mean(fitness))
            }
            return (CoV(meanFitnesses) < fitness_plateau_tolerance)
        case "cov": //CoV of fitness for last generation
            return (CoV(fitnesses[len(fitnesses)-1]) < fitness_plateau_tolerance)
        default:
            panic("Invalid plateau mode, must be {cov_mean|std_dev}")
        }
}
// RunSimulation() runs a genetic algorithm and returns the final generation
// input: simulation parameters as commented
// output: last (most fit) generation (list of members)
func RunSimulation(lower int,
                   upper int,
                   size int,
                   maxIterations int,
                   targetFile string,
                   mutation_rate float64,
                   indel_rate float64,
                   top_sequence_percent float64,
                   fitness_mode string,
                   fitness_plateau_tolerance float64,
                   plateau_gens int) Population {
    target := ReadTargetFromFasta(targetFile)
    currentGen := InitializeGeneration(size,lower,upper,target)
    bar := pb.StartNew(maxIterations).Prefix("Generations:")
    var generationFitnesses [][]float64 //list of fitness values for all solutions for each generation
    for gen := 0; gen < maxIterations; gen++ {//terminate regardless after maxIterations
        //keep only last plateau_gens generational fitnesses stored
        generationFitnesses = generationFitnesses[Max(0,len(generationFitnesses)-plateau_gens):len(generationFitnesses)]
        generationFitnesses = append(generationFitnesses,currentGen.FitnessList())
        //check if average fitness has plateaued
        if FitnessPlateau(fitness_mode,generationFitnesses,fitness_plateau_tolerance,) {
            bar.Finish()
            fmt.Println("Reached Fitness Plateau at generation ",gen)
            return currentGen //if plateau, no improvements from continnuing simulation, finish
        }
        currentGen = BreedNewGeneration(currentGen,target,mutation_rate,indel_rate,top_sequence_percent)
        bar.Increment()
    }//never reached fitness plateau
    bar.Finish()
    fmt.Println("Reached Max Iterations ",maxIterations)
    return currentGen
}

