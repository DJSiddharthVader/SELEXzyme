package main

import(
    "fmt"
)

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
// Complementarity() returns BLAST score of sequence to target, higher is better for fitness
// output: float64
func (s Member) Complementarity(target Member) float64 {
    return float64(len(s.seq)+len(target.seq))
}
// ScoreFitness() asseses the total fitness of a sequence
// output: fitness score of s.seq
func (s Member) ScoreFitness() float64 {
    return float64(len(s.seq))
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
*/
func newline() {
    fmt.Println()
}
