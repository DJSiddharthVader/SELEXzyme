package main

import(
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
