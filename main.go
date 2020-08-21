package main

import (
	"fmt"
	"math"
	"regexp"
	"strings"
)

const k = .5

var rgxAlphaNum = regexp.MustCompile("[^a-z0-9]")

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
	SpamCount      int             // count of spam messages
	HamCount       int             // count of ham messages
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

// train is a method on classifier that "trains" the Bayes classifier
func (cl *classifier) fit(msgs []message) {
	// Tally up the number of spam and ham messages in the set of message objects
	for _, msg := range msgs {
		if msg.isSpam {
			cl.SpamCount = cl.SpamCount + 1
		} else {
			cl.HamCount = cl.HamCount + 1
		}

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
	}

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
func (cl *classifier) predict(txt string) float64 {
	var ok bool
	var lnPSpam, lnPHam float64
	var pSpam, pHam float64

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

	return pSpam / (pSpam + pHam)
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

func main() {
	myMessages := []message{
		{"spam rules", true},
		{"ham rules", false},
		{"hello ham", false},
	}

	myModel := newClassifier(k)
	myModel.fit(myMessages)

	myPred := myModel.predict("hello spam")

	fmt.Println(myPred)
	foo()
}
