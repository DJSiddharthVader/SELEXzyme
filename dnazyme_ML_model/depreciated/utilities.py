from sklearn.feature_extraction.text import CountVectorizer

KMER = 6
TRAIN_RATIO = 0.75
VALIDATION_RATIO = 0.15
TEST_RATIO = 0.10
VECTORIZER = CountVectorizer(ngram_range=(4, 4))  # arbitrary, can be optimized


def get_kmers(seq, size):
    return [seq[x:x+size].upper() for x in range(len(seq) - size + 1)]


def get_words(seq, size):
    return ' '.join(get_kmers(seq, size))
