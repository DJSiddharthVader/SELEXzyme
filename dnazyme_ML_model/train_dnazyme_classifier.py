import pickle
import pandas as pd
from sklearn import svm
import utilities as util
from sklearn.model_selection import train_test_split
import sklearn.metrics as metrics


def evaluate(y_test, y_pred):
    accuracy = metrics.accuracy_score(y_test, y_pred)
    precision = metrics.precision_score(y_test, y_pred, average='weighted')
    recall = metrics.recall_score(y_test, y_pred, average='weighted')
    f1 = metrics.f1_score(y_test, y_pred, average='weighted')
    confusion_matrix = pd.crosstab(pd.Series(y_test, name='Actual'),
                                   pd.Series(y_pred, name='Predicted'))
    return confusion_matrix, accuracy, precision, recall, f1


def prep_data(seq_tsv):
    df = pd.read_csv(seq_tsv, sep='\t')
    df = df[['Sequence', 'Label_Binary', 'Label']]
    indices = list(range(0, 1000)) + list(range(20000, 21000))
    df = df.iloc[indices, :]
    # create kmers for each sequence
    df['Words'] = df['Sequence'].apply(lambda x: util.get_words(x, util.KMER))
    X = util.VECTORIZER.fit_transform(df['Words']).toarray()
    y = df['Label_Binary']
    return X, y, df


def main(seq_tsv, outfile):
    test_size = 1-util.TRAIN_RATIO
    validation_size = util.TEST_RATIO/(util.TEST_RATIO + util.VALIDATION_RATIO)
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
    confusion, accuracy, precision, recall, f1 = evaluate(y_test, y_pred_SVM)
    print(confusion)
    print(accuracy, precision, recall, f1)
    print(dict(zip(classifier_SVM.classes_,
                   classifier_SVM.decision_function(X_test[0]))))
    pickle.dump(classifier_SVM, open(outfile, 'wb'))
    return None


if __name__ == '__main__':
    seq_tsv = '../data/All_Seqeunces.tsv'
    outfile = './dnazyme_classifier_SVM.sav'
    main(seq_tsv, outfile)
