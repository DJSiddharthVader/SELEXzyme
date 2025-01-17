import os
import sys
import random
import pandas as pd
from tqdm import tqdm
from Bio import SeqIO

"""
Prepares training data for the dnazyme learning model
parses fasta sequences of DNAzymes, aptamers and promoter sequences
into a csv with labels as input for training a model to score if a
target sequences is a DNAzyme.
Also include some randomly generated sequences as negative examples
in addition to the aptamers and promoters.
Outputs a tsv file with the sequences,  fasta header info and labels
for training.
"""

# Files
data_dir = os.path.abspath('../data')
dnazyme_fasta_filepath = os.path.join(data_dir, 'NCBI_DNAzymes.fasta')
aptamer_fasta_filepath = os.path.join(data_dir, 'NCBI_Aptamers.fasta')
promoter_fasta_filepath = os.path.join(data_dir, 'NCBI_Promoters.fasta')

DNA_ALPHABET = 'ACGT'


def GetGCContentDistribution(sequences):
    """ GetGCContentDistribution() get the distribution of GC content
    from a collection of sequences
    input: sequence iterable
    output: scipy.stats distribution object
    """
    sequences = [seq.upper() for seq in sequences]
    return [(seq.count('G')+seq.count('C'))/len(seq) for seq in sequences]


def MakeRandomSequence(length, gcContent):
    """
    MakeRandomSequence() produces a random sequence of length L from
    input: length of sequence,  percent of bases that are G or C
    output: random sequence (each base is picked randomly but with constraints
    """
    gc_prob = gcContent/2
    at_prob = (1 - gc_prob)/2
    base_probs = [at_prob, gc_prob, gc_prob, at_prob]
    return ''.join(random.choices(DNA_ALPHABET, k=length, weights=base_probs))


def MakeRandomSequenceDf(n, lower, upper, gcContentDistribution,
                         label, binary_label):
    """
    MakeRandomSequenceDf() produces a pandas df of N random sequences
    of legnth lower to upper
    input: n size of dataframe,  lower and upper are lower/upper
           bounds of the sequence length
    output: pandas dataframe of random sequences
    """
    rows = []
    for i, gc_content in tqdm(enumerate(gcContentDistribution)):
        length = random.choice(range(lower, upper))  # random length
        sequence = MakeRandomSequence(length, gc_content)
        row = {'GI': 'RSID_{}'.format(i),  # for pd.concat()
               'Identifier_Type': 'RSID',
               'Identifier_Value': 'RSID_{}'.format(i),
               'Description': 'Random Sequence {}'.format(i),
               'Sequence': sequence}
        rows.append(row)
    df = pd.DataFrame(rows)
    df['Label'] = label  # classification labels
    df['Label_Binary'] = binary_label  # binary,  DNAzyme or not DNAzyme
    return df


def ExtractHeader(header):
    """
    ExtractHeader() parses a Fasta header into a dict
    headers are of the form:
    >gi|GeneInfo ID|IDENTIFIER_TYPE|IDENTIFIER_VALUE|$description
    Example header: >gi|551461007|dbj|DI203508.1| KR 1020120021470-A/11:
     Method for Detecting Nucleic Acid Using DNAzyme Having Peroxidase Activity
    What Identifiers can be is found here:
        https://en.wikipedia.org/wiki/FASTA_format#NCBI_identifiers
    input: header string
    output: dictionary
    """
    fields = header.split('|')
    header_dict = {'GI': fields[1],
                   'Identifier_Type': fields[2].upper(),
                   'Identifier_Value': fields[3],
                   'Description': fields[4]}
    return header_dict


def FastaToDf(fasta_filepath, label, binary_label):
    """
    FastaToDf() produces a pandas dataframe with seqeunces in on column and
                several metadata columns,  extracted from the fasta header
    input: fasta_filepath
    output: pandas dataframe containnig all the fasta sequences
    """
    record_dicts = []  # each entry in the fasta will be a dict
    records = list(SeqIO.parse(fasta_filepath, 'fasta'))
    for record in tqdm(records):
        record_dict = ExtractHeader(record.description)
        record_dict.update({'Sequence': str(record.seq)})
        record_dicts.append(record_dict)
    df = pd.DataFrame(record_dicts)
    df['Label'] = label  # classification labels
    df['Label_Binary'] = binary_label  # binary,  DNAzyme or not DNAzyme
    return df


def MakeData():
    """
    MakeData() produces entire training set for the dnazyme evaluation model
    input: none
    output: dataframe of all sequences and labels
    """
    # Make DNAzyme dataframe
    dnazyme_df = FastaToDf(dnazyme_fasta_filepath, 'DNAzyme', 'DNAzyme')
    # Make Aptamer dataframe
    aptamer_df = FastaToDf(aptamer_fasta_filepath, 'Aptamer', 'Not_DNAzyme')
    # Make Promoter dataframe
    promoter_df = FastaToDf(promoter_fasta_filepath, 'Promoter', 'Not_DNAzyme')
    # Make Random dataframe
    total_sequences = sum([dnazyme_df.shape[0],
                           aptamer_df.shape[0],
                           promoter_df.shape[0]])
    gcDistribution = GetGCContentDistribution(dnazyme_df['Sequence'])
    minDNAzyme = min([len(x) for x in dnazyme_df.Sequence])
    maxDNAzyme = max([len(x) for x in dnazyme_df.Sequence])
    random_df = MakeRandomSequenceDf(int(total_sequences/3),
                                     minDNAzyme, maxDNAzyme, gcDistribution,
                                     'Random', 'Not_DNAzyme')
    # Combine all dataframes together
    all_df = pd.concat([dnazyme_df, aptamer_df, promoter_df, random_df])
    return all_df


if __name__ == '__main__':
    if len(sys.argv) < 2:
        outfilepath = os.path.join(data_dir, 'All_Seqeunces.tsv')
    else:
        outfilepath = sys.argv[1]
    df = MakeData()
    df.to_csv(outfilepath,  sep='\t',  index=False)
    """
    df['GC'] = [(seq.count('G')+seq.count('C'))/len(seq)
                for seq in df.Sequence]
    print(df[df.Label == 'DNAzyme'][['Label', 'GC']].describe())
    print(df[df.Label == 'Random'][['Label', 'GC']].describe())
    """

""" DEPRECIATED
def MakeRandomSequence(length, gcContentProportion):
    MakeRandomSequence() produces a random sequence of length L from
    input: length of sequence,  percent of bases that are G or C
    output: random sequence (each base is picked randomly but with constraints
    bases = []
    for i in range(0, length):
        gcProbability = random.uniform(0, 1)
        if gcProbability < gcContentProportion:
            base = random.choice(['G', 'C'])
        else:
            base = random.choice(['A', 'T'])
        bases.append(base)
    return ''.join(bases)
"""
