# from sklearn import CountVectorizer
# from sklearn.naive_bayes import MultinomialNB
import pandas as pd

def read_wrangle_csv(fname):
    df = pd.read_csv(fname,
    names=["Category", "Message"])
    print(df.describe())

read_wrangle_csv("../data/standardSpamData.csv")