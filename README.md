SELEXzyme: Generating DNAzymes for Target Sequences using a Genetic Algorithm
=============================================================================

## Demo Idea
- pick an miRNA
- evolve dnazymes
- plot various sequence stats against fitness (GC content, length, etc.)
- compare alignment to known in vivo verified dnazymes
- estimate fitness for known seqs and comapre to generations

## Project Overview

## Logistics

### Commands and Running
This program used python to train the machine learning model for the cost function and implements the genetic algorithm in golang.
The likelihood that a sequence is a DNAzyme is produced by calling `./src/dnazyme_ML_model/evaluate_sequence.py $MODEL_FILE $SEQ` where `$MODEL_FILE` contains the machine learning model parameters.
Running the program will be as follows

`./genetic_algorithm $target.fasta $output_file $num_gens $mutation_rate $population_size $model_file $output_format`

The parameters are
 - `$target.fasta` fasta file with 1 entry contain the sequence you want a DNAzyme to target
 - `$output_file` output file containing a population of dnazyme sequences
   - can output tsv or fasta file, automatically detected from extension
   - fasta includes sequence label and fitness in headers
   - tsv include columns:
     - __seqeunce:__ the DNAzyme sequence
     - __id:__ unique identifier for this sequence
     - __hairpins:__ number of hairpins found
       __dnazyme_score:__ the score output from our dnazyme model described [here](#machine-learning-dnazyme-classification-model)
     - __fitness:__ total fitness score as described in [above](#fitness-function)
   - `$num_gens` maximum number of generations to simulate if fitness does not plateau before
   - `$mutation_rate` mutation rate for new sequences described [here](#mutation)
   - `$population_size` total number of sequences in the population
   - `$model_file` file with parameters for model making the dnazyme prediction

you can also run `./genetic_algorithm -h` to see a list of arguments and defaults

### External Dependencies

#### Python
A conda environment yml is included for the python dependencies
- get from conda output

#### Golang
- github.com/biogo/biogo
- github.com/cheggaaa/pb

## Genetic Algorithm

### Breeding
The crux of a genetic algorithm is breeding a new population of solutions from a current population of solutions to increase the overall fitness.
This way after a number of iterations we will have a much fitter population than the initial population.
In our case a solution is an actual DNA string representing a DNAzyme that will specifically bind our user-specified target sequence.
So given a population of solutions $P_i$ we first take the fittest members of the population $F \subset P_i$ (default top 20%) and use them to breed a new population $P_{i+1}$.
Breeding here essentially means crossover and mutate to create a new solution.
So we randomly select two solutions $a,b \in F$ (with replacement)  and create a new solution
$$c = mutate(crossover(a,b))$$
We continue generating $c$'s until our new population is the same size as our current population ($|P_{i+1}| = |P_i|$)
$$P_{i+1} = F \cup \{c_i \forall i=0; i < |P_i|-|F|; i++ \}$$
Note that $F$ is passed on to the new population as these are still our best solutions so far.

#### Crossover
Crossover is a key element of genetic algorithms.
Crossover helps ensure that new solutions are created during each generation while retaining the best features from the current set of solutions.
It is done by taking parts of two different solutions and gluing them together to form a new solution.
In this case our solutions are actual genetic sequences so the process is fairly trivial.
First we select a random index $i$ and then merge two sequences as such
$$seq_3 = seq_1[0:i] + seq_2[i:]$$
taking the first part of $seq_1$ and the second part of $seq_2$.
This will always result in a new $seq_3$ of length $Max(len(seq_1,seq_2))$.

#### Mutation
Again mutation is implemented to introduce more variation into the solution population, specifically to avoid getting stuck in local optima.
For mutation every base in a sequence is mutated with some probability $\mu$ (default 0.005).
This is done to more closely mimic actual models of DNA mutation, as opposed to mutating a set percentage of bases for each sequence.
Further there is a chance (default 0.1) that a given mutation can result in an insertion or deletion (each 0.5 probability).
In a deletion that base is deleted from the solution, in an insertion a new random base is added after the current base, which is left unchanged.
Mutations must change the base to a new base, so you cannot have a $T \to T$ mutation.


### Halting
At some point the program must halt and cease to produce new generations of solutions.
This is done in two cases; reaching the max number of iterations (default 500) or no increases in population fitness.
This is an optimization algorithm so eventually it will reach a fitness optimum and no longer be able to improve its solutions.
To asses this we first calculate the average fitness of all solutions for each generation separately.
Next we check if the coefficient of variation (std.dev / mean) of the previous $g$ generations (default 5) is less than some threshold $\theta$ (default 0.25)
This means that the average fitness has not changed much in the past $g$ generations and thus further generations will not change the average fitness enough.
This is the default option but you can also specify to just consider if the CoV of fitness from the current generation is below a threshold.


### Fitness Function

#### Hairpins

#### Complementarity To Target
The complementarity to the target sequence is also considered when assessing sequence fitness, specifically
$$\frac{Smith-Waterman(sequence,target)}{len(target)}$$
is what is calculated.
The Smith-Waterman local alignment score is calculated between the current DNAzyme and the target sequence.
We use Smith-Waterman because only part of the DNAzyme needs to match the target and the DNAzymes are likely larger than the target regions.
This is the raw score, counting mismatched and gaps according to the scoring matrix defined [here](./genetic_algorithm/data.go).
This score is then divided by the length of the target, ensuring that the maximum score is always 1, for a perfect matching sequence.

#### Catalytic Activity
See [here](#machine-learning-dnazyme-classification-model)

## Machine Learning DNAzyme Classification Model

### Data Collection
Data is required for training the machine learning model to asses how likely a given DNA sequence is a DNAzyme.
There are relatively few known DNAzyme sequence so we also want to include negative examples, which include random DNA sequences and aptamer like sequences that can bind proteins but have no catalytic activity.
Known DNAzymes, DNA APtamers and promoters were downloaded from the [NCBI Nucleotide](https://www.ncbi.nlm.nih.gov/nuccore/), by searching for "DNAzyme","Aptamer" and "promoter" respectively and filtering for "genomic DNA/RNA".
Promoters were constrained to be between 50-250 base pairs long to trim the results and more closely resemble Aptamers and DNAzymes
Random sequences between 50-300bp were also generated as negative training examples.
A tsv including all the sequences and labels is provided in `./data/all_sequences.csv`, which contains
 - 9999 Aptamers
 - 5581 DNAzymes
 - 4537 Promoters
 - 6706 Random sequences

An example file is viewable [here](./data/example_training_data.tsv)
(Note values the `Identifier_Type` column are described [here](https://en.wikipedia.org/wiki/FASTA_format#NCBI_identifiers).

### Training/Algorithms

### Model Evaluation

