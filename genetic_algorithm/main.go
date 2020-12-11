package main

import(
    "os"
    "fmt"
    "flag"
    "strings"
    "math/rand"
)

// Between() check if target is between max and min inclusive
// input: float to check, floor ,ceiling
// output: bool, true if in range
func Between(target float64, min float64, max float64) bool {
    if target >= min && target <= max {
        return true
    } else {
        return false
    }
}
// CheckPrams() check if the supplied parameters are valid and panics otherwise
func CheckParams(lower int,
                 upper int,
                 size int,
                 maxIterations int,
                 targetFile string,
                 mutation_rate float64,
                 top_sequence_percent float64,
                 fitness_plateau_tolerance float64,
                 fitness_plateau_generations int,
                 outputfile string) {
    /* Parameter Restrictions
    - All numerical values must be positive
    - mutation rate must be in [0,1]
    - top_Seqs must be in [0,1]
    - lower < upper
    - plateau generations < maxIterations
    - target file must be a fasta file (handeled in target reading func)
    - target must be a DNA string (A,G,C,T only) (handeled by biogo)
    - outputfile must contain a valid extension (checked in write results)
    DEP
    - len(target) < upper, otherwise could not cover target sequence
    */
    //Check positiives
    numericalParams := []interface{}{lower,upper,size,maxIterations,fitness_plateau_tolerance,fitness_plateau_generations}
    for i,param := range numericalParams {
        switch param.(type) {
            case int:
                if param.(int) < 0 {
                    panic(fmt.Sprint("Param",i,"invalid, must be >= 0"))
                }
            case float64:
                if param.(float64) < 0 {
                    panic(fmt.Sprint("Param",i,"invalid, must be >= 0"))
                }
            default:
                panic("Invalid numeric type, must be int or float64")
        }
    }
    //Check specific restrictions
    switch false {//check if each statement is false (indicates error)
        case lower <= upper:
            panic("lower >= upper")
        case Between(mutation_rate,0,1):
            panic("mutation rate must be in [0,1]")
        case Between(top_sequence_percent,0,1):
            panic("top_sequence_Percent must be in [0,1]")
        case fitness_plateau_generations < maxIterations:
            panic("generatoins to consider for fitness plateau must be < maxIterations")
    }
}

func main() {
    rand.Seed(9) //for testing
    //Parse Arguments
    //Initialization Params
    lower := flag.Int("lower",10,"minimum length of initial sequences")
    upper := flag.Int("upper",100,"maximum length of initial seuqences")
    size := flag.Int("size",1000,"number of sequences in the populations")
    maxIterations := flag.Int("maxIters",30,"max generations to simulate")
    targetFastaFile := flag.String("target","target.fna","target sequence for generated dnazymes to catalyze")

    //Simulation params
    mutation_rate := flag.Float64("mutation",0.005,"mutation rate for sequences, in [0,1]")
    indel_rate := flag.Float64("indel",0.1,"probability for mutation being an indel, in [0,1]")
    top_sequence_percent := flag.Float64("top_seqs",0.2,"percentage of sequences to use for breeding, in [0,1]")
    model_file := flag.String("model","../dnazyme_ML_model/dnazyme_SGD_Classifier_v1.pickle","model used for DNAzyme evaluation (pickle of sklearn model")
    // minimum_hairpin_length := flag.Int("hairpin_len",4,"minimum size for a sequence to be considered pallindromic")

    //Termination Params
    eval := flag.String("eval","","Only evaluates the fitness of sequences in fasta passed")
    fitness_plateau_mode := flag.String("plateau","cov_mean","criteria for deciding on fitness plateau, one of {cov_mean|cov}")
    fitness_plateau_tolerance := flag.Float64("plateau_tol",0.005,"maximum CoV of previous generations of fitness when deciding on plateau")
    fitness_plateau_generations := flag.Int("plateau_gens",5,"number of generations to consider for evaluating fitness plateau")
    outputfile := flag.String("output","dnazymes.fna","output file name for final set of dnazymes, must have extension {.tsv|.fna}")

    flag.Parse()
    CheckParams(*lower,
                *upper,
                *size,
                *maxIterations,
                *targetFastaFile,
                *mutation_rate,
                *top_sequence_percent,
                *fitness_plateau_tolerance,
                *fitness_plateau_generations,
                *outputfile)

    //Run simulation
    if len(*eval) != 0 { //only evaluate fitness of input fasta
        pop := FastaToPopulation(*eval)
        target := ReadTargetFromFasta(*targetFastaFile)
        pop.ScoreFitness(target, *model_file)
        outfile := fmt.Sprintf("%s_fitness.fna",strings.Replace(*eval,".fna","",-1))
        pop.WriteToFasta(outfile)
        fmt.Println("Scored file written to ",outfile)
        os.Exit(0) //exit without simulating
    }
    lastGen := RunSimulation(*lower,
                             *upper,
                             *size,
                             *maxIterations,
                             *targetFastaFile,
                             *mutation_rate,
                             *indel_rate,
                             *top_sequence_percent,
                             *model_file,
                             *fitness_plateau_mode,
                             *fitness_plateau_tolerance,
                             *fitness_plateau_generations)
    fmt.Println("Final Generation Fitness Summary")
    lastGen.Summarize()
    lastGen.WriteResults(*outputfile)
}
