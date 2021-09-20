package main

import(
    "os"
    "io"
    "fmt"
    "sort"
    "math"
    "math/rand"
    "strings"
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
// Max() returns maximum of 2 ints
func Max(x,y int) int {
    if x > y {
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
// StdDev() Stanrad deviation of n numbers
func StdDev(n []float64) float64 {
    mean := Mean(n)
    total := 0.0
    for _,f := range n {
        total += math.Pow(f-mean,2)
    }
    return math.Sqrt(total/float64(len(n)-1))
}
// CoV() coefficient of variance of a list of normalized numbers
func CoV(n []float64) float64 {
    return math.Abs(StdDev(n)/Mean(n))
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

// FitnessList() returns a slice of all fitness values for pop members
func (pop Population) FitnessList() []float64 {
    fitnesses := make([]float64,len(pop))
    for i,member := range pop {
        fitnesses[i] = member.fitness
    }
    return fitnesses
}
// SortByFitness() returns a copy of pop sorted ascending order of fitness
// so the fittest sequences are at the end of the sorted slice
func (pop Population) SortByFitness() Population {
    sort.Slice(pop,func(i,j int) bool { return pop[i].fitness < pop[j].fitness })
    return pop
}
// MeanFitness() get mean fitness of every Sequence in a population
func (pop Population) MeanFitness() float64 {
    return Mean(pop.FitnessList())
}
// Summarize() prints some summary statistics for the fitnesses of a population
func (pop Population) Summarize() {
    fitnesses := pop.SortByFitness().FitnessList()
    fmt.Println("-------------------------------------")
    fmt.Println("Min..............",fitnesses[0]) //sorted
    fmt.Println("25% Quart        ",fitnesses[len(fitnesses)/4])
    fmt.Println("Mean.............",Mean(fitnesses))
    fmt.Println("75% Quart        ",fitnesses[3*len(fitnesses)/4])
    fmt.Println("Max..............",fitnesses[len(fitnesses)-1])
    fmt.Println("Std. Dev         ",StdDev(fitnesses))
    fmt.Println("CoV..............",CoV(fitnesses))
    fmt.Println("-------------------------------------")
}

// ReadTargetFromFasta() takes a fasta file and returns the first entry as a *linear.Seq object
// input: fasta file name
// output: biogo *linear.Seq object, can be aligned
func ReadTargetFromFasta(fastafilename string) *linear.Seq {
    fastaFile, err := os.Open(fastafilename)
    if err != nil { panic(err) }
    defer fastaFile.Close()

    template := linear.NewSeq("",alphabet.Letters{},alphabet.DNA)
    reader := fasta.NewReader(fastaFile,template)
    targetSeq, err := reader.Read() //read first seq in file
    if err != nil { panic(err) }
    var sequence string //DNA sequence
    for i:=0; i<targetSeq.Len(); i++ {//add one letter at a time
        sequence += strings.ToUpper(string(targetSeq.At(i).L))
    }
    target := linear.Seq{Seq:[]alphabet.Letter(sequence)}
    target.Alpha = ALPHABET
    return &target
}
// FastaToPopulation() reads a fasta file into a Population object
// input: fasta file name
// output: Population ([]Member)
func FastaToPopulation(fastafilename string) Population {
    fastaFile, err := os.Open(fastafilename)
    if err != nil { panic(err) }
    defer fastaFile.Close()

    template := linear.NewSeq("",alphabet.Letters{},alphabet.DNA)
    reader := fasta.NewReader(fastaFile,template)
    var pop Population
    for {
        seq,err := reader.Read()
        if err == io.EOF {
            break
        }
        var sequence string //DNA sequence
        for i:=0; i<seq.Len(); i++ {//add one letter at a time
            sequence += strings.ToUpper(string(seq.At(i).L))
        }
        member := Member{seq:sequence,
                         header:seq.CloneAnnotation().ID}
        pop = append(pop,member)
    }
    return pop
}
// ConvertToSeqObject() converts a Member to a seq.Sequence (biogo object) for file writing
// input: member object
// output: biogo.linear.Seq object that can be writen easily
func (s Member) ConvertToSeqObject() seq.Sequence {
    var label string
    if s.header == "" {
        label = fmt.Sprintf("Sequence_%v | Fitness:%v",s.label,s.fitness,)
    } else {
        label = fmt.Sprintf("%v | Fitness:%v",s.header,s.fitness,)
    }
    return linear.NewSeq(label,[]alphabet.Letter(s.seq),alphabet.DNA)
}

// WriteToFasta () write every member of the population into a fasta file
// input: population, list of members
func (pop Population) WriteToFasta(filename string) {
    //convert to seq.Sequence object
    outfile,err := os.Create(filename)
    if err != nil {
        panic(err)
    }
    defer outfile.Close()
    writer := fasta.NewWriter(outfile,80) //width 80
    for _,member := range pop {
        writer.Write(member.ConvertToSeqObject())
    }
}
// WriteToTSV() write every member of the population into a tsv file
// input: population, list of members
func (pop Population) WriteToTSV(filename string) {
    //convert to seq.Sequence object
    outfile,err := os.Create(filename)
    if err != nil {
        panic(err)
    }
    defer outfile.Close()
    line := fmt.Sprint("Index\tSeqLabel\tFitness\tSequence\n")
    outfile.WriteString(line)
    for i,member := range pop {
        line = fmt.Sprintf("%d\tSequence_%d\t%f\t%s\n",i,member.label,member.fitness,member.seq)
        outfile.WriteString(line)
    }
}
// WriteResults () write every member of the population into a file, either tsv or fasta
// output file may conatin < len(pop) entries Because it removes duplicates
// input: population, list of members
// output: no return, write file to filename
func (pop Population) WriteResults(filename string) {
    splits := strings.Split(filename,".") //last element of split is the extension
    extension := splits[len(splits)-1] //last element of split is the extension
    switch extension {
        case "fna":
            pop.WriteToFasta(filename)
        case "tsv":
            pop.WriteToTSV(filename)
        default:
            panic(fmt.Sprintf("Invalid outfile format, must be {fna|tsv}"))
    }
}
