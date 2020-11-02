package main

import(
    "os"
    "fmt"
    "sort"
    "math"
    "math/rand"
    "github.com/biogo/biogo/alphabet"
    "github.com/biogo/biogo/seq"
    "github.com/biogo/biogo/seq/linear"
    "github.com/biogo/biogo/io/seqio/fasta"
)

// Min() returns minimum of 2 ints
func Min(x,y int) int {
    if x < y {
        return x
    } else {
        return y
    }
}
// Mean() mean of a list of numbers
func Mean(n []float64) float64 {
    total := 0.0
    for _,f := range n {
        total += f
    }
    return total/float64(len(n))
}
// CoV() coefficient of variance of a list of normalized numbers
func CoV(n []float64) float64 {
    mean := Mean(n)
    total := 0.0
    for _,f := range n {
        total += math.Pow(f-mean,2)
    }
    stdDev := math.Sqrt(total/float64(len(n)-1))
    return math.Abs(stdDev/mean)
}

// RandomIntBetween() returns a random in between 2 other ints
// input: lower and upper bounds
// output: random int between lower and upper
// from https://flaviocopes.com/go-random/
func RandomIntBetween(lower,upper int) int {
    if lower >= upper {
        panic("lower must be strictly smaller than upper")
    }
    return lower + rand.Intn(upper-lower)
}
// PickRandomBase() picks a random DNA base
func PickRandomBase() string {
    return string(DNA_ALPHABET[rand.Intn(len(DNA_ALPHABET))])
}
// PickDifferentRandomBase() picks a random DNA base that
// is different from the base you pass as an argument
func PickDifferentRandomBase(base rune) string {
    // var DNA_ALPHABET = [4]rune{'A','C','G','T'}
    var baseIndex int
    switch base {
        case 'A':// A is at position 0, avoid it
            baseIndex = []int{1,2,3}[rand.Intn(3)]
        case 'C':// C is at position 1, avoid it
            baseIndex = []int{0,2,3}[rand.Intn(3)]
        case 'G':// G is at position 2, avoid it
            baseIndex = []int{0,1,3}[rand.Intn(3)]
        case 'T':// T is at position 3, avoid it
            baseIndex = []int{0,1,2}[rand.Intn(3)]
    }
    return string(DNA_ALPHABET[baseIndex])
}

// SortByFitness() returns a copy of pop sorted ascending order of fitness
// so the fittest sequences are at the end of the sorted slice
func (pop Population) SortByFitness() Population {
    sort.Slice(pop,func(i,j int) bool { return pop[i].fitness < pop[j].fitness })
    return pop
}
// MeanFitness() get mean fitness of every Sequence in a population
func (pop Population) MeanFitness() float64 {
    fitnesses := make([]float64,len(pop))
    for i,Seq := range pop {
        fitnesses[i] = Seq.fitness
    }
    return Mean(fitnesses)
}

// ConvertToSeqObject() converts a Member to a seq.Sequence (biogo object) for file writing
// input: member object
// output: biogo.linear.Seq object that can be writen easily
func (s Member) ConvertToSeqObject() seq.Sequence {
    label := fmt.Sprintf("Sequence_%v | Fitness:%v",s.label,s.fitness,)
    alphaSeq := alphabet.BytesToLetters([]byte(s.seq))
    return linear.NewSeq(label,alphaSeq,alphabet.DNA)
}
// WriteToFasta () write every member of the population into a file, either tsv or fasta
// output file may conatin < len(pop) entries vecause it removes duplicates
// input: population, list of members
// output: no return, write file to filename
func (pop Population) WriteToFasta(filename string) {
    //convert to seq.Sequence object
    outfile,err := os.Create(filename)
    if err != nil {
        panic(err)
    }
    defer outfile.Close()
    writer := fasta.NewWriter(outfile,80) //width 80
    fmt.Println(len(pop))
    for i,member := range pop {
        writer.Write(member.ConvertToSeqObject())
        fmt.Println(i)
    }
}
/*
*/
