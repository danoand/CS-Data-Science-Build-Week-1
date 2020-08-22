package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/danoand/utils"
	"github.com/gin-gonic/gin"
)

const k = .5
const stdMdlDataFname = "data/standardSpamData.csv"

var err error
var rgxAlphaNum = regexp.MustCompile("[^a-z0-9]")
var stdModel *classifier

// handlerPred handles a request to generate a spam prediction for a
//     given string of text
func handlerPred(c *gin.Context) {
	var err error
	var bbytes []byte
	var prms = make(map[string]string)
	var rMap = make(map[string]interface{})

	// Has the standard model been trained?
	if !stdModel.HasBeenFit {
		// model has not been flagged as trained
		rMap["msg"] = "model has not yet been trained"
		c.JSON(http.StatusInternalServerError, rMap)
		return
	}

	// Grab the request body data
	bbytes, err = c.GetRawData()
	if err != nil {
		// error fetching the request body
		log.Printf("ERROR: %v - error fetching the request body. See: %v\n",
			utils.FileLine(),
			err)
		rMap["msg"] = "error fetching the request body"
		c.JSON(http.StatusInternalServerError, rMap)
		return
	}

	// Parse the request body
	err = utils.FromJSONBytes(bbytes, &prms)
	if err != nil {
		// error parsing the request body data
		log.Printf("ERROR: %v - error parsing the request body data. See: %v\n",
			utils.FileLine(),
			err)
		rMap["msg"] = "error parsing the request body data"
		c.JSON(http.StatusInternalServerError, rMap)
		return
	}

	// Missing text value?
	if len(prms["text"]) == 0 {
		// error missing text value
		log.Printf("ERROR: %v - error missing text value\n",
			utils.FileLine())
		rMap["msg"] = "error missing text value"
		c.JSON(http.StatusBadRequest, rMap)
		return
	}

	// Invoke the prediction
	pred := stdModel.predict(prms["text"])
	rMap["msg"] = "prediction of spam"
	rMap["prediction"] = pred
	rMap["text"] = prms["text"]
	c.JSON(http.StatusOK, rMap)
}

func main() {
	stdModel, err = fitSpamModel2Data("data/standardSpamData.csv")
	if err != nil {
		log.Fatalf("ERROR: %v - error fitting the standard model. See: %v\n",
			utils.FileLine(),
			err)
	}
	foo()

	prdStr := "07732584351 - Rodger Burns - MSG = We tried to call you re your reply to our sms for a free nokia mobile + free camcorder. Please call now 08000930705 for delivery tomorrow"
	myPred := stdModel.predict(prdStr)

	fmt.Println(myPred)
	foo()

	// Set up a gin web server
	r := gin.Default()
	r.POST("/predict", handlerPred)

	// Start the web server
	port := "localhost:8090"
	log.Printf("INFO: %v - running web server on port: %v\n",
		utils.FileLine(),
		port)
	r.Run("localhost:8090")
}
