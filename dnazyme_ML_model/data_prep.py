import os
import sys
import random
import pandas as pd
from tqdm import tqdm
from Bio import SeqIO

"""
Prepares training data for the dnazyme learning model
parses fasta sequences of DNAzymes,aptamers and promoter sequences
into a csv with labels as input for training a model to score if a
target sequences is a DNAzyme.
Also include some randomly generated sequences as negative examples
in addition to the aptamers and promoters.
Outputs a tsv file with the sequences, fasta header info and labels
for training.
"""

# Files
data_dir = os.path.abspath('../data')
dnazyme_fasta_filepath = os.path.join(data_dir,'NCBI_DNAzymes.fasta')
aptamer_fasta_filepath = os.path.join(data_dir,'NCBI_Aptamers.fasta')
promoter_fasta_filepath = os.path.join(data_dir,'NCBI_Promoters.fasta')

# Constants
DNA_ALPHABET = 'ACGT'

# MakeRandomSequence() produces a random sequence of length L from
def MakeRandomSequence(length):
    return ''.join(random.choice(DNA_ALPHABET) for i in range(0,length))

# MakeRandomSequenceDf() produces a pandas df of N random sequences of legnth lower to upper
# input: n size of dataframe, lower and upper are lower/upper bounds of the sequence length
# output: pandas dataframe of random sequences
def MakeRandomSequenceDf(n,lower,upper,label,binary_label):
    rows = []
    for i in tqdm(range(n)):
        length = random.choice(range(lower,upper)) #random length between lower and upper
        row = {'GI':'RSID_{}'.format(i),    # not particularly important, just used as unique identifiers and allows for pd.concat()
               'Identifier_Type':'RSID',
               'Identifier_Value':'RSID_{}'.format(i),
               'Description':'Random Sequence {}'.format(i),
               'Sequence':MakeRandomSequence(length)}
        rows.append(row)
    df = pd.DataFrame(rows)
    df['Label'] = label #classification labels
    df['Label_Binary'] = binary_label #binary, DNAzyme or not DNAzyme
    return df

# ExtractHeader() parses a Fasta header into a dict
# headers are of the form: >gi|GeneInfo ID|IDENTIFIER_TYPE|IDENTIFIER_VALUE|$description of the sequence
# Example header: >gi|551461007|dbj|DI203508.1| KR 1020120021470-A/11: Method for Detecting Nucleic Acid Using DNAzyme Having Peroxidase Activity
# What Identifiers can be is found here: https://en.wikipedia.org/wiki/FASTA_format#NCBI_identifiers
# input: header string
# output: dictionary
def ExtractHeader(header):
    fields = header.split('|')
    header_dict = {'GI':fields[1],
                   'Identifier_Type':fields[2].upper(),
                   'Identifier_Value':fields[3],
                   'Description':fields[4]}
    return header_dict

# FastaToDf() produces a pandas dataframe with seqeunces in on column and several metadata columns, extracted from the fasta header
# input: fasta_filepath
# output: pandas dataframe containnig all the fasta sequences
def FastaToDf(fasta_filepath,label,binary_label):
    record_dicts = [] #each entry in the fasta will be a dict, parsed using SeqIO.Parse
    records = list(SeqIO.parse(fasta_filepath,'fasta'))
    for record in tqdm(records):
        record_dict = ExtractHeader(record.description)
        record_dict.update({'Sequence':str(record.seq)})
        record_dicts.append(record_dict)
    df = pd.DataFrame(record_dicts)
    df['Label'] = label #classification labels
    df['Label_Binary'] = binary_label #binary, DNAzyme or not DNAzyme
    return df

# MakeData() produces entire training set for the dnazyme evaluation model
# input: none
# output: dataframe of all sequences and labels
def MakeData():
    # Make DNAzyme dataframe
    dnazyme_df = FastaToDf(dnazyme_fasta_filepath,'DNAzyme','DNAzyme')
    # Make Aptamer dataframe
    aptamer_df = FastaToDf(aptamer_fasta_filepath,'Aptamer','Not_DNAzyme')
    # Make Promoter dataframe
    promoter_df = FastaToDf(promoter_fasta_filepath,'Promoter','Not_DNAzyme')
    # Make Random dataframe
    total_sequences = sum([dnazyme_df.shape[0],aptamer_df.shape[0],promoter_df.shape[0]])
    random_df = MakeRandomSequenceDf(int(total_sequences/3),50,300,'Random','Not_DNAzyme')
    # Combine all dataframes together
    all_df = pd.concat([dnazyme_df,aptamer_df,promoter_df,random_df])
    return all_df


if __name__ == '__main__':
    if len(sys.argv) < 2:
        outfilepath = os.path.join(data_dir,'All_Seqeunces.tsv')
    else:
        outfilepath = sys.argv[1]
    df = MakeData()
    df.to_csv(outfilepath,sep='\t',index=False) #use tabs because some descriptions contain commas
