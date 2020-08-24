from sklearn.feature_extraction.text import CountVectorizer
from sklearn.model_selection import train_test_split
from sklearn.naive_bayes import MultinomialNB
import pandas as pd

# class2int returns 1 for a 'spam' class; 0 otherwise
def class2int(val):
    if val == "spam":
        return 1

    return 0

# read_wrangle_csv reads in the model csv file and applies
#    basic data wrangling
def read_wrangle_csv(fname):
    df = pd.read_csv(fname, names=["Category", "Message"])
    df['is_spam'] = df['Category'].apply(class2int)
    df.drop(df.index[0], inplace=True)

    return df

# gen_model generates a vectorizer and model object and fits 
#   a model given the dataset processed at startup
def gen_model(DF):
    print(DF.describe())
    X = DF['Message']
    y = DF['is_spam']

    X_train, X_test, y_train, y_test = train_test_split(X, y)

    # Stand up a vectorizer object
    vctrizr = CountVectorizer()
    # Generate the word token counts
    tkn_counts = vctrizr.fit_transform(X_train.values)

    # Stand up a classifier object
    clssifier = MultinomialNB()
    # Fit the classifier model
    model_fit = clssifier.fit(tkn_counts, y_train.values)

    # Return the fit model and the vectorizer for use downstream
    return model_fit, vctrizr

# predict predicts a spam classification using the passed model, vectorizer, and 
#    prediction source data
def mdl_pred(mdl, vtzr, val):
    # Vectorize the inbound data
    val_counts = vtzr.transform(val)
    # Generate a spam prediction
    pred       = mdl.predict(val_counts)

    return pred

# --- MAIN ---
# Create and clean up the underlying model dataset
df = read_wrangle_csv("../data/standardSpamData.csv")
# Generate and fit a Bayesian Classifier model
model, vctrizr  = gen_model(df)
# Generate a model prediction
pred = mdl_pred(
    model,
    vctrizr,
    ["Even my brother is not like to speak with me. They treat me like aids patent."])
print("pred is: ", pred)

