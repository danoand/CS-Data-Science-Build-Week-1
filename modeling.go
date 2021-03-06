package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"strings"

	"github.com/danoand/utils"
)

// msgStore serves as a mini db or store of messages for use downstream
var msgStore = make(map[string][]message)

// message defines a message instance to be trained (and used for predictions)
type message struct {
	Text   string
	isSpam bool
}

// classifier defines a model classifier
type classifier struct {
	K              float64         // smoothing factor
	Tokens         map[string]bool // map: set of classifier tokens
	TokenSpamCount map[string]int  // map: count of spam tokens
	TokenHamCount  map[string]int  // map: count of ham tokens
	TrnObsCount    int             // count of training observations
	SpamCount      int             // count of spam messages
	HamCount       int             // count of ham messages
	HasBeenFit     bool            // indicator that model training has been exected
}

// newClassifier creates a new classifier object to be used downstream
func newClassifier(k float64) *classifier {
	var tmpCls classifier

	tmpCls.K = k
	tmpCls.Tokens = make(map[string]bool)
	tmpCls.TokenSpamCount = make(map[string]int)
	tmpCls.TokenHamCount = make(map[string]int)

	return &tmpCls
}

// fit is a method on classifier that "fits" the Bayes classifier
func (cl *classifier) fit(msgs []message) {
	cl.TrnObsCount = 0

	// Tally up the number of spam and ham messages in the set of message objects
	for _, msg := range msgs {

		// Increment the spam and ham token (word) counts
		_, mapTkns := tokenize(msg.Text)
		for tkn := range mapTkns {
			// Add the token to the classifier list
			cl.Tokens[tkn] = true

			// Increment token spam count
			if msg.isSpam {
				cl.TokenSpamCount[tkn] = cl.TokenSpamCount[tkn] + 1
			} else {
				cl.TokenHamCount[tkn] = cl.TokenHamCount[tkn] + 1
			}
		}

		cl.TrnObsCount = cl.TrnObsCount + 1
	}

	// Flag that the model has been "fit"/trained
	cl.HasBeenFit = true

	return
}

// calcBayesProb is a method on classifier that calculates the probability of a token
//    for a given class (spam or ham) probabilities for a given token
func (cl *classifier) calcBayesProb(tkn string) (float64, float64) {
	spam := cl.TokenSpamCount[tkn]
	ham := cl.TokenHamCount[tkn]

	probSpam := (float64(spam) + cl.K) / (float64(cl.SpamCount) + 2.0*cl.K)
	probHam := (float64(ham) + cl.K) / (float64(cl.HamCount) + 2.0*cl.K)

	return probSpam, probHam
}

// predict is a method on classifier that predicts the probability that a message is "spam"
func (cl *classifier) predict(txt string) (float64, float64, float64) {
	var (
		ok              bool
		lnPSpam, lnPHam float64
		pSpam, pHam     float64
	)

	// tokenize the inbound string of text
	_, mapTkns := tokenize(txt)

	// Iterate through each classifier token
	for tkn := range cl.Tokens {
		// Calculate the spam and ham probablities of the current token
		pSpam, pHam = cl.calcBayesProb(tkn)

		// Is this token in the the text
		_, ok = mapTkns[tkn]
		if ok {
			// current token is in this message
			lnPSpam = math.Log(pSpam) + lnPSpam
			lnPHam = math.Log(pHam) + lnPHam
		} else {
			// current token is NOT in this message
			lnPSpam = math.Log(1.0-pSpam) + lnPSpam
			lnPHam = math.Log(1.0-pHam) + lnPHam
		}
	}

	// Calculate the overall probability of spam and ham
	pSpam = math.Exp(lnPSpam)
	pHam = math.Exp(lnPHam)
	numr := pSpam
	dnom := pSpam + pHam

	return (pSpam / (pSpam + pHam)), numr, dnom
}

// tokenize returns a slice of word tokens extracted from a string of text
func tokenize(txt string) ([]string, map[string]bool) {
	var (
		rtArr []string
		rtMap = make(map[string]bool)
	)

	txt = strings.ToLower(txt)
	arrWords := strings.Split(txt, " ")

	// Iterate through the slice of words and extract tokens
	for _, wrd := range arrWords {
		// remove non alphaneric characters
		tmpWrd := rgxAlphaNum.ReplaceAllString(wrd, "")

		// add the word to the unique word map
		rtMap[tmpWrd] = true
	}

	// Construct a slice of unique word tokens
	for k := range rtMap {
		rtArr = append(rtArr, k)
	}

	return rtArr, rtMap
}

func foo() {
	return
}

// fitSpamModel2Data fits the "house" or standard classifier model
//    1. loads csv file from disk
//    2. creates a model for use by the frontend
//    3. persists the model for reuse (future)
func fitSpamModel2Data(fname string) (*classifier, map[string][]message, error) {
	var (
		err            error
		ctrSpm, ctrHam int
		fcsv           *os.File
		flines         [][]string
		tmpModel       *classifier
		msgs           []message
		msgStr         = make(map[string][]message)
	)

	// Initialize the message "store"
	msgStr["spam"] = []message{}
	msgStr["ham"] = []message{}

	log.Printf("INFO: %v - fitting the standard spam model for use downstream\n",
		utils.FileLine())

	// Validate the filenmame parameter
	if len(fname) == 0 {
		log.Printf("ERROR: %v - error - missing data filename\n",
			utils.FileLine())

		return tmpModel, msgStr, fmt.Errorf("missing filename")
	}

	// Create a new classifier/model object
	tmpModel = newClassifier(k)

	// Open the data file containing the underlying data for the standard model
	fcsv, err = os.Open(stdMdlDataFname)
	if err != nil {
		log.Printf("ERROR: %v - error opening the model data file: %v. See: %v\n",
			utils.FileLine(),
			stdMdlDataFname,
			err)

		return tmpModel, msgStr, err
	}

	// Read the csv file
	flines, err = csv.NewReader(fcsv).ReadAll()
	if err != nil {
		log.Printf("ERROR: %v - error reading the csv data file. See: %v\n",
			utils.FileLine(),
			err)

		return tmpModel, msgStr, err
	}
	log.Printf("INFO: %v - read in the standard model data with %v rows\n",
		utils.FileLine(),
		len(flines))

	// Iterate through the csv data and generate an array of message objects
	//    skip the first row (assume it's a header row)
	for i := 1; i < len(flines); i++ {
		var tmpMsg = message{}
		var tmpLne = flines[i]

		// Ignore invalid lines
		if len(tmpLne[0]) == 0 || len(tmpLne[1]) == 0 {
			continue
		}

		// Assign the text value to the "message" object
		tmpMsg.Text = tmpLne[1]

		if tmpLne[0] == "spam" {
			tmpMsg.isSpam = true
			ctrSpm = ctrSpm + 1
			msgStr["spam"] = append(msgStr["spam"], tmpMsg)
		} else {
			ctrHam = ctrHam + 1
			msgStr["ham"] = append(msgStr["ham"], tmpMsg)
		}

		// append the new message object to our messages array
		msgs = append(msgs, tmpMsg)
	}

	// Assign the count value of spam and ham messages respectively
	tmpModel.SpamCount = ctrSpm
	tmpModel.HamCount = ctrHam

	// Fit the temporary model using the default dataset
	tmpModel.fit(msgs)

	log.Printf("INFO: %v - fit model with %v message datapoints [%v, %v]\n",
		utils.FileLine(),
		len(msgs),
		ctrSpm,
		ctrHam)

	return tmpModel, msgStr, err
}
