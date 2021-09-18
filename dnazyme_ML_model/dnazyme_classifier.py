import sys
import pickle
from Bio import SeqIO
from sklearn.feature_extraction.text import HashingVectorizer

kmer_len = 6  # size of kmers
n_features = 4**kmer_len+1  # number of unique possible kmers, avoids collision
VECTORIZER = HashingVectorizer(n_features=n_features)  # for str to kmer hash


def fasta_to_list(fasta_file):
    """ Convert a fasta file into a list of string objects
        input: fasta file name
        output: list of strings, 1 per fasta entry
    """
    return [str(rec.seq) for rec in SeqIO.parse(fasta_file, "fasta")]


# turn a string into a list of non-overlapping kmers (1 space delimited string)
def get_kmers(seq, size=kmer_len):
    return ' '.join([seq[x:x+size].upper() for x in range(len(seq)-size+1)])


# Vectorize list of kmers to format suitable for model prediction
def seq_to_vector(fasta_file):
    # list of kmers of size kmer_len per seq
    seq_list = fasta_to_list(fasta_file)
    kmer_lists = [get_kmers(seq) for seq in seq_list]
    return VECTORIZER.transform(kmer_lists).toarray()


def main(fasta_file, model_file):
    # load trained model
    model = pickle.load(open(model_file, 'rb'))
    # prep sequence to evaluate
    X = seq_to_vector(fasta_file)
    # prob that label of seq is 1 (DNAzyme) according to model
    predictions = model.predict_proba(X)
    return [x[1] for x in predictions]


if __name__ == '__main__':
    fasta_file = sys.argv[1]
    model_file = sys.argv[2]
    output = main(fasta_file, model_file)
    print(' '.join(["%.9f" % x for x in output]))
