import sys
import numpy as np
from sklearn import svm


def get_kmers(seq, size):
    return [seq[x:x+size].upper() for x in range(len(seq) - size + 1)]


def get_words(seq, size):
    return ' '.join(get_kmers(seq, size))


def main(sequence, kmer_size):
    words = get_words(sequence, kmer_size)
    classifier_SVM = svm.SVC(probability=False, C=0.5)  # expect noise
    # stackoverflow.com/questions/49507066/predict-probabilities-using-svm
    pred = classifier_SVM.decision_function(words)[0]
    prob = np.exp(pred)/np.sum(np.exp(pred))  # softmax
    return prob


if __name__ == '__main__':
    sequence = sys.argv[1]
    kmer_size = 6
    main(sequence, kmer_size)
