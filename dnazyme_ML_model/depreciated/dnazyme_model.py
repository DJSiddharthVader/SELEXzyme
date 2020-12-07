import sys
import pickle
import numpy as np
import pandas as pd
from sklearn import svm
import sklearn.metrics as metrics
from sklearn.model_selection import train_test_split
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


def prep_data(seq_tsv):
    df = pd.read_csv(seq_tsv, sep='\t')
    df = df[['Sequence', 'Label_Binary', 'Label']]
    indices = list(range(0, 1000)) + list(range(20000, 21000))
    df = df.iloc[indices, :]
    # create kmers for each sequence
    df['Words'] = df['Sequence'].apply(lambda x: get_words(x, KMER))
    X = VECTORIZER.fit_transform(df['Words']).toarray()
    y = df['Label_Binary']
    return X, y, df


def evaluate_model(y_test, y_pred):
    accuracy = metrics.accuracy_score(y_test, y_pred)
    precision = metrics.precision_score(y_test, y_pred, average='weighted')
    recall = metrics.recall_score(y_test, y_pred, average='weighted')
    confusion_matrix = pd.crosstab(pd.Series(y_test, name='Actual'),
                                   pd.Series(y_pred, name='Predicted'))
    return confusion_matrix, accuracy, precision, recall


def train_model(seq_tsv, model_file):
    test_size = 1-TRAIN_RATIO
    validation_size = TEST_RATIO/(TEST_RATIO + VALIDATION_RATIO)
    # Load Features (bag of words) labels (binary) and entire df
    X, y, df = prep_data(seq_tsv, test_size, validation_size)
    # Splitting data
    X_train, X_test, y_train, y_test = train_test_split(X, y,
                                                        test_size=test_size)
    X_val, X_test, y_val, y_test = train_test_split(X_test, y_test,
                                                    test_size=validation_size)
    # initialize and train SVM classifier
    classifier_SVM = svm.SVC(probability=True, C=0.5)  # expect noise
    classifier_SVM.fit(X_train, y_train)
    y_pred_SVM = classifier_SVM.predict(X_test)
    # evaluate model accuracy
    confusion, accuracy, precision, recall = evaluate_model(y_test, y_pred_SVM)
    print(confusion)
    print(accuracy, precision, recall)
    print(dict(zip(classifier_SVM.classes_,
                   classifier_SVM.decision_function(X_test[0]))))
    pickle.dump(classifier_SVM, open(model_file, 'wb'))
    return None


def classify(sequence, model_file, kmer_size):
    # load trained model
    model = pickle.load(open(model_file, 'wb'))
    # prep sequence to evaluate
    words = get_words(sequence, KMER)  # convert seq to bagofwords
    X = VECTORIZER.fit_transform([words]).toarray()  # length 1
    # "normalize" SVM decision score to probability
    # stackoverflow.com/questions/49507066/predict-probabilities-using-svm
    pred = model.decision_function(X)[0]
    prob = np.exp(pred)/np.sum(np.exp(pred))  # softmax
    return prob


if __name__ == '__main__':
    seq_tsv = '../data/All_Seqeunces.tsv'
    model_file = './dnazyme_classifier_SVM.sav'
    if sys.argv[1] == 'train':
        train_model(seq_tsv, model_file)
    elif sys.argv[1] == 'classify':
        classify(sys.argv[2], model_file)
    else:
        raise ValueError("Invalid mode")
