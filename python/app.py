# Modeling imports  
from sklearn.feature_extraction.text import CountVectorizer
from sklearn.model_selection import train_test_split
from sklearn.naive_bayes import MultinomialNB
import pandas as pd

# Flask imports
from flask import Flask
from flask import request
from flask import jsonify

app = Flask(__name__)

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

# --- FLASK APP ---
# /status can be used to verify that the Flask app is running and responding
@app.route('/status')
def test():
  global df
  num_rows = len(df.index)
  ret_dict = {}
  ret_dict["msg"] = "app.py is up and running"
  ret_dict["note"] = "the number of rows in the dataset is: " + str(num_rows)

  return jsonify(ret_dict)

# /predict accepts a string of text and returns a spam prediction value
@app.route('/predict', methods = ['POST'])
def predict():
  global model
  global vctrizr

  # Get the request's body as json
  jsn      = request.get_json(force=True)
  ret_dict = {}

  # Validate the inbound json body parameters
  tmp_string = jsn["text"]
  if len(tmp_string) == 0:
    # missing text value
    ret_dict["msg"] = "missing text value"
    return jsonify(ret_dict), 400

  # Process the inbound text and generate a list of tokens
  tmp_pred = mdl_pred(model, vctrizr, [tmp_string])

  # Construct the return object
  return_object = {}
  return_object["msg"]              = "Scikit Learn spam prediction"
  return_object["content"]          = dict(
    prediction=str(float(tmp_pred[0])*100.0)+"%",
    havepred=True
    )
  return jsonify(return_object)

# Start the flask app
if __name__ == '__main__':
  app_domain = 'localhost'
  app_port   = 8091
  print(f'INFO: starting web api server on {app_domain}:{app_port}')
  app.run(host='localhost', port=app_port)


