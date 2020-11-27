package main

import(
    "flag"
    "math/rand"
    "github.com/biogo/biogo/alphabet"
    "github.com/biogo/biogo/seq/linear"
)
/* Parameter Restrictions
*/

func main() {
    rand.Seed(9) //for testing
    //Parse Arguments
    //Initialization Params
    lower := flag.Int("lower",50,"minimum length of initial sequences")
    upper := flag.Int("upper",500,"maximum length of initial seuqences")
    size := flag.Int("size",500,"number of sequences in the populations")
    maxIterations := flag.Int("maxIters",30,"max generations to simulate")
    //Target Sequence object for DNAzymes to bind to
    dnazymeTarget := flag.String("target","","target sequence for generated dnazymes to catalyze")
    target := &linear.Seq{Seq:alphabet.BytesToLetters([]byte(*dnazymeTarget))}
    target.Alpha = ALPHABET

    //Simulation params
    mutation_rate := flag.Float64("mutation",0.005,"mutation rate for sequences")
    top_sequence_percent := flag.Float64("top_seqs",0.2,"percentage of sequences to use for breeding")
    // minimum_hairpin_length := flag.Int("hairpin_len",4,"minimum size for a sequence to be considered pallindromic")

    //Termination Params
    fitness_plateau_tolerance := flag.Float64("plateau_tol",0.25,"maximum CoV of previous generations of fitness when deciding on plateau")
    fitness_plateau_generations := flag.Int("plateau_gens",5,"number of generations to consider for evaluating fitness plateau")
    outputfile := flag.String("output","dnazymes.fna","output file name for final set of dnazymes")

    flag.Parse()
    //Run simulation
    lastGen := RunSimulation(*lower,
                             *upper,
                             *size,
                             *maxIterations,
                             target,
                             *mutation_rate,
                             *top_sequence_percent,
                             *fitness_plateau_tolerance,
                             *fitness_plateau_generations)
    lastGen.WriteResults(*outputfile)
}
