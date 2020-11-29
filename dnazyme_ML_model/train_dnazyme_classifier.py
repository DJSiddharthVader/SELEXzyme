import pandas as pd
from sklearn import svm
from sklearn.model_selection import train_test_split
from sklearn.feature_extraction.text import CountVectorizer
import sklearn.metrics as metrics


def get_kmers(seq, size):
    return [seq[x:x+size].upper() for x in range(len(seq) - size + 1)]


def get_words(seq, size):
    return ' '.join(get_kmers(seq, size))


def evaluate(y_test, y_pred):
    accuracy = metrics.accuracy_score(y_test, y_pred)
    precision = metrics.precision_score(y_test, y_pred, average='weighted')
    recall = metrics.recall_score(y_test, y_pred, average='weighted')
    f1 = metrics.f1_score(y_test, y_pred, average='weighted')
    confusion_matrix = pd.crosstab(pd.Series(y_test, name='Actual'),
                                   pd.Series(y_pred, name='Predicted'))
    return confusion_matrix, accuracy, precision, recall, f1


def main(kmer_size, train_ratio, test_ratio, validation_ratio):
    test_size = 1-train_ratio
    validation_size = test_ratio/(test_ratio + validation_ratio)

    data_file = '../data/All_Seqeunces.tsv'
    df = pd.read_csv(data_file, sep='\t')
    df = df[['Sequence', 'Label_Binary', 'Label']]
    indices = list(range(0, 1000)) + list(range(20000, 21000))
    df = df.iloc[indices, :]
    # create kmers for each sequence
    df['Words'] = df['Sequence'].apply(lambda x: get_words(x, kmer_size))
    cv = CountVectorizer(ngram_range=(4, 4))  # arbitrary, can be optimized
    X = cv.fit_transform(df['Words']).toarray()
    y = df['Label_Binary']
    X_train, X_test, y_train, y_test = train_test_split(X, y,
                                                        test_size=test_size)
    X_val, X_test, y_val, y_test = train_test_split(X_test, y_test,
                                                    test_size=validation_size)

    classifier_SVM = svm.SVC(probability=True, C=0.5)  # expect noise
    classifier_SVM.fit(X_train, y_train)
    y_pred_SVM = classifier_SVM.predict(X_test)
    confusion, accuracy, precision, recall, f1 = evaluate(y_test, y_pred_SVM)
    print(confusion)
    print(accuracy, precision, recall, f1)
    print(dict(zip(classifier_SVM.classes_,
                   classifier_SVM.decision_function(X_test[0]))))


if __name__ == '__main__':
    kmer_size = 6
    train_ratio = 0.75
    validation_ratio = 0.15
    test_ratio = 0.10
    main(kmer_size, train_ratio, test_ratio, validation_ratio)
