package main

import(
    "fmt"
    "os/exec"
    "strconv"
    "strings"
    "github.com/biogo/biogo/alphabet"
    "github.com/biogo/biogo/seq/linear"
)

const classifer_script = "../dnazyme_ML_model/dnazyme_classifier.py"
const tmp_fasta = "/tmp/generation_tmp.fna"
const python_exe = "/home/sidreed/anaconda3/envs/selexzyme/bin/python3"

// Complementarity() returns BLAST score of sequence to target, higher is better for fitness
// output: float64, local (SW) alignment score of the 2 sequences
// thoguh target is an argument it will be constant through out the simulation
// as it will always be the user supplied target sequence
func (s Member) Complementarity(target *linear.Seq) float64 {
    seq := &linear.Seq{Seq:alphabet.BytesToLetters([]byte(s.seq))}
    seq.Alpha = ALPHABET //set alphabet, required by biogo
    aln, err := SW_MATRIX.Align(seq,target)
    if err != nil { panic(err) }
    swScore := aln[0].(Scorer).Score()
    return float64(swScore)/float64(Min(len(s.seq),len(target.Seq)))
}
// CallDNAzymeModel() call a machine learning model to estimate
// the likelihood  that this sequence is a DNAzyme
func (pop Population) CallDNAzymeModel(model_file string) []float64 {
    pop.WriteToFasta(tmp_fasta)
    cmd := exec.Command(python_exe, classifer_script, tmp_fasta, model_file)
    out, err := cmd.Output()
    if err != nil {
        panic(err)
    }
    output := strings.Split(strings.TrimSuffix(string(out),"\n")," ")
    predictions := make([]float64,len(pop))
    for i, prediction := range output {
        predictions[i], _ = strconv.ParseFloat(prediction, 64) // return model probability as a float
    }
    return predictions
}
// ScoreFitness() asseses the total fitness every sequence in a population
// output: no return, fitness is assigned for every seq inplace
func (pop Population) ScoreFitness(target *linear.Seq, model_file string) {
    predictions := pop.CallDNAzymeModel(model_file)
    for i,member := range pop {
        similarity := member.Complementarity(target)
        dnazymeness := predictions[i]
        pop[i].fitness = (similarity*0.4+dnazymeness*0.6)/2
    }
}

/*DEPRECIATED
// MeltingTemperature() get the melting temperature of the sequence, returns abs difference to target
// output: float64
func (s Sequence) MeltingTemperature(targetTemp float64) float64 {
    return float64(len(s.seq))
}
// HasHairpins() check if the sequence contains a palindrome that will cause hairpin binding in vivo
// output: int, number of hairpins
func (s Sequence) HasHairpins() int {
    hairpins := 0
    totalHalfwayPoint := int(len(s.seq)/2)
    var halfwayPoint int
    var firstHalf,secondHalfReverse string
    for windowSize := MINIMUM_HAIRPIN_LENGTH; windowSize < totalHalfwayPoint; windowSize++ {
        // check for pallendroms of length ranging from 2*minimumLength bp to half of the entire sequence
        // 0.5 times length of the pallendrome
        for windowStart := 0; windowStart <= len(s.seq)-windowSize; windowStart++{
            // check every position if it is a pallendroms of length windowSize
            halfwayPoint = windowStart+windowSize // halfway point of the potential pallendrome
            firstHalf = s.seq[windowStart:halfwayPoint] // first half of the potential pallendrome
            secondHalfReverse = s.ReverseComplement()
            if firstHalf == secondHalfReverse[halfwayPoint:halfwayPoint+windowSize] {
                // even length palindrome ACGG CCGT
                hairpins += 1
            } else if firstHalf == secondHalfReverse[halfwayPoint:halfwayPoint+windowSize+1] {
                // odd length palindrome ACGG T CCGT
                hairpins += 1
            }
        }
    }
    return hairpins
}
// ScoreFitness() asseses the total fitness of a sequence
// output: fitness score of s.seq
func (s Member) ScoreFitness(target *linear.Seq, model_file string) float64 {
    similarity := s.Complementarity(target)
    dnazymeness := s.CallDNAzymeModel(model_file)
    return (similarity*0.4+dnazymeness*0.6)/2
}
// ReverseComplement() reverses a DNA string and takes the complement
// output: string
func (s Member) ReverseComplement() string {
    var result string
    for _,v := range s.seq {
        result = string(DNA_COMPLEMENTS[v]) + result
    }
    return result
}
// HasHairpins() check if the sequence contains a palindrome that will cause hairpin binding in vivo
// output: int, number of hairpins
func (s Member) HasHairpins() int {
    hairpins := 0
    totalHalfwayPoint := int(len(s.seq)/2)
    var halfwayPoint int
    var firstHalf,secondHalfReverse string
    for windowSize := minimum_hairpin_length; windowSize < totalHalfwayPoint; windowSize++ {
        // check for pallendroms of length ranging from 2*minimumLength bp to half of the entire sequence
        // 0.5 times length of the pallendrome
        for windowStart := 0; windowStart <= len(s.seq)-windowSize; windowStart++{
            // check every position if it is a pallendroms of length windowSize
            halfwayPoint = windowStart+windowSize // halfway point of the potential pallendrome
            firstHalf = s.seq[windowStart:halfwayPoint] // first half of the potential pallendrome
            secondHalfReverse = s.ReverseComplement()
            if firstHalf == secondHalfReverse[halfwayPoint:halfwayPoint+windowSize] {
                // even length palindrome ACGG CCGT
                hairpins += 1
            } else if firstHalf == secondHalfReverse[halfwayPoint:halfwayPoint+windowSize+1] {
                // odd length palindrome ACGG T CCGT
                hairpins += 1
            }
        }
    }
    return hairpins
}
*/
func newline() { fmt.Println() }
