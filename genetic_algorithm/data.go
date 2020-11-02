package main

//Simulation parameters
const OPTIMAL_MELTING_TEMP = 37 //human body temperature in celcius
const MINIMUM_HAIRPIN_LENGTH = 4 //only count hairpins of length 4
const MUTATION_RATE = 0.005 //mutation rate for sequences
const TOP_SEQUENCE_PERCENT = 0.2 //percentage of sequences to use for breeding
const FITNESS_PLATEAU_TOLERANCE= 0.25 //maximum std dev. of average seq fitness of past generations until stopping
const FITNESS_PLATEAU_GENERATIONS = 5 //number of generations for which fintess must have plateaued to halt simulation

//Constants
var DNA_ALPHABET = [4]rune{'A','C','G','T'}
var DNA_COMPLEMENTS = map[rune]rune{'A':'T','C':'G','G':'C','T':'A'}

//Data Types
type Member struct {
    seq string
    fitness float64
    label int
}
type Population []Member
