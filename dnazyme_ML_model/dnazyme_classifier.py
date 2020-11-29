import sys
import pickle
import numpy as np
import utilities as util


def main(sequence, model_file, kmer_size):
    # load trained model
    model = pickle.load(open(model_file, 'wb'))
    # prep sequence to evaluate
    words = util.get_words(sequence, util.KMER)  # convert seq to bagofwords
    X = util.VECTORIZER.fit_transform([words]).toarray()  # length 1
    # "normalize" SVM decision score to probability
    # stackoverflow.com/questions/49507066/predict-probabilities-using-svm
    pred = model.decision_function(X)[0]
    prob = np.exp(pred)/np.sum(np.exp(pred))  # softmax
    return prob


if __name__ == '__main__':
    sequence = sys.argv[1]
    model_file = sys.argv[2]
    main(sequence, model_file)
