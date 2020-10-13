SELEXzyme: Generating DNAzymes for Target Sequences using a Genetic Algorithm
=============================================================================

## Project Overview

## Genetic Algorithm

### Breeding

#### Crossover

#### Mutation

#### Selection

### Fitness Function

#### Hairpins

#### Melting Temperature

#### Complementarity To Target

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

## Logistics

### Commands and Running
This program used python to train the machine learning model for the cost function and implements the genetic algorithm in golang.
The likelihood that a sequence is a DNAzyme is produced by calling `./src/dnazyme_ML_model/evaluate_sequence.py $MODEL_FILE $SEQ` where `$MODEL_FILE` contains the machine learning model parameters.
Running the program will be as follows

`./src/genetic_algorithm/genetic_algorithm $target.fasta $output_file $num_gens $mutation_rate $population_size $model_file $output_format`

The parameters are
 - `$target.fasta` fasta file with 1 entry containg the sequence you want a DNAzyme to target
 - `$output_file` csv file containing a population of solutions with sequence info in the header
   - csv include columns:
     - __seqeunce:__ the DNAzyme sequence
     - __id:__ unique identifier for this sequence
     - __hairpins:__ number of hairpins found
     - __melting_temp:__ melting temperature
       __dnazyme_score:__ the score output from our dnazyme model described [here](#machine-learning-dnazyme-classification-model)
     - __fitness:__ total fitness score as described in [above](#fitness-function)
   - `$output_format` either csv or fasta, fasta will contain csv info in the header of each sequence
   - `$num_gens` maximum number of generations to simulate if fitness does not plateau before
   - `$mutation_rate` mutation rate for new sequences described [here](#mutation)
   - `$population_size` total number of sequences in the population
   - `$model_file` file with paramaters for model making the dnazyme prediction

### Dependencies
A conda environment yml and a requirements.txt are included for the python dependencies

#### Programs
 - blastn

#### Python
 - Biopython
 - Pandas
 - tqdm

#### Golang
 -

