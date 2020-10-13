package main

import(
    "fmt"
    "os/exec"
    "strings"
)

// ReverseComplement() reverses a DNA string and takes the complement
// output: string
func (s Sequence) ReverseComplement() string {
    var COMPLEMENT_MAP = map[rune]rune{'A':'T','T':'A','G':'C','C':'G'} //complementarty base pairs
    var result string
    for _,v := range s.seq {
        result = string(COMPLEMENT_MAP[v]) + result
    }
    return result
}
// HasHairpins() check if the sequence contains a palindrome that will cause hairpin binding in vivo
// output: int, number of hairpins
func (s Sequence) HasHairpins(minimumnLengthProportion float64) int {
    hairpins := 0
    minimumLength := int(float64(len(s.seq))*minimumnLengthProportion)
    totalHalfwayPoint := int(len(s.seq)/2)
    var halfwayPoint int
    var firstHalf,secondHalfReverse string
    for windowSize := minimumLength; windowSize < totalHalfwayPoint; windowSize++ {
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
// Complementarity() returns BLAST score of sequence to target, higher is better for fitness
// output: float64
func (s Sequence) Complementarity(target Sequence) float64{
    BLASTN_PATH := "/usr/bin/blastn"
    commandString := strings.Split(fmt.Sprintf("echo -e \">query\n%s\" >| query.fasta; echo -e \">suject\n%s\" >| subject.fasta; %s -query query.fasta -subject subject.fasta",s.seq,target.seq,BLASTN_PATH)," ")
    command :=  exec.Command(commandString[0],commandString[1:]...)
    fmt.Println(command.Run())
    result,err := command.Output()
    fmt.Println(result)
    fmt.Println(err)
    return float64(len(s.seq)+len(target.seq))
}
/*
// MeltingTemperature() get the melting temperature of the sequence, returns abs difference to target
// output: float64
func (s Sequence) MeltingTemperature(targetTemp float64) float64{
    return float64(len(s.seq))
}
// Score() asseses the total fitness of a sequence
// output: no output, just updates the fitness score
func (s Sequence) Score(target string) Sequence {
    var t Sequence
    t.fitness = float64(len(s.seq)+len(target))
}
*/
