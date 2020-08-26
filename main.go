package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"regexp"

	"github.com/danoand/utils"
	"github.com/gin-gonic/gin"
)

const k = .5
const stdMdlDataFname = "data/standardSpamData.csv"
const URLPyMLSvc = "http://localhost:8091/predict"

var err error
var rgxAlphaNum = regexp.MustCompile("[^a-z0-9]")
var stdModel *classifier

// getModelInfo returns information on the currently applicable classifier
func getModelInfo(c *gin.Context) {
	var respMap = make(map[string]interface{})
	var rMap = make(map[string]interface{})

	rMap["hasbeentrained"] = stdModel.HasBeenFit //
	rMap["numobs"] = stdModel.TrnObsCount        //
	rMap["numtkns"] = len(stdModel.Tokens)       //
	rMap["numspam"] = stdModel.SpamCount         // number of spam messages
	rMap["numham"] = stdModel.HamCount           // number of ham messages

	respMap["content"] = rMap

	c.JSON(http.StatusOK, respMap)
}

// hdlrPred handles a request to generate a spam prediction for a
//     given string of text
func hdlrPred(c *gin.Context) {
	var err error
	var bbytes []byte
	var prms = make(map[string]string)
	var respMap = make(map[string]interface{})
	var cMap = make(map[string]interface{})

	// Has the standard model been trained?
	if !stdModel.HasBeenFit {
		// model has not been flagged as trained
		respMap["msg"] = "model has not yet been trained"
		cMap["havepred"] = false
		respMap["content"] = cMap
		c.JSON(http.StatusInternalServerError, respMap)
		return
	}

	// Grab the request body data
	bbytes, err = c.GetRawData()
	if err != nil {
		// error fetching the request body
		log.Printf("ERROR: %v - error fetching the request body. See: %v\n",
			utils.FileLine(),
			err)
		respMap["msg"] = "an error occurred, please try again"
		cMap["havepred"] = false
		respMap["content"] = cMap
		c.JSON(http.StatusInternalServerError, respMap)
		return
	}

	// Parse the request body
	err = utils.FromJSONBytes(bbytes, &prms)
	if err != nil {
		// error parsing the request body data
		log.Printf("ERROR: %v - error parsing the request body data. See: %v\n",
			utils.FileLine(),
			err)
		respMap["msg"] = "an error occurred, please try again"
		cMap["havepred"] = false
		respMap["content"] = cMap
		c.JSON(http.StatusInternalServerError, respMap)
		return
	}

	// Missing text value?
	if len(prms["text"]) == 0 {
		// error missing text value
		log.Printf("ERROR: %v - error missing text value\n",
			utils.FileLine())
		respMap["msg"] = "error missing text value"
		cMap["havepred"] = false
		respMap["content"] = cMap
		c.JSON(http.StatusBadRequest, respMap)
		return
	}

	// Invoke the prediction
	pred := stdModel.predict(prms["text"])
	respMap["msg"] = "a spam prediction"
	cMap["havepred"] = true
	cMap["prediction"] = fmt.Sprintf("%.2f", pred*100.0)
	respMap["content"] = cMap
	c.JSON(http.StatusOK, respMap)
}

// hdlrPyPred handles a request to generate a prediction using Python libraries
func hdlrPyPred(c *gin.Context) {
	var (
		err             error
		bbytes, reqBody []byte
		rspBody         []byte
		resp            *http.Response
		prms            = make(map[string]string)
		respMap         = make(map[string]interface{})
		respSvc         = make(map[string]interface{})
		cMap            = make(map[string]interface{})
	)

	// Grab the request body data
	bbytes, err = c.GetRawData()
	if err != nil {
		// error fetching the request body
		log.Printf("ERROR: %v - error fetching the request body. See: %v\n",
			utils.FileLine(),
			err)
		respMap["msg"] = "an error occurred, please try again"
		cMap["havepred"] = false
		respMap["content"] = cMap
		c.JSON(http.StatusInternalServerError, respMap)
		return
	}

	// Parse the request body
	err = utils.FromJSONBytes(bbytes, &prms)
	if err != nil {
		// error parsing the request body data
		log.Printf("ERROR: %v - error parsing the request body data. See: %v\n",
			utils.FileLine(),
			err)
		respMap["msg"] = "an error occurred, please try again"
		cMap["havepred"] = false
		respMap["content"] = cMap
		c.JSON(http.StatusInternalServerError, respMap)
		return
	}

	// Missing text value?
	if len(prms["text"]) == 0 {
		// error missing text value
		log.Printf("ERROR: %v - error missing text value\n",
			utils.FileLine())
		respMap["msg"] = "error missing text value"
		cMap["havepred"] = false
		respMap["content"] = cMap
		c.JSON(http.StatusBadRequest, respMap)
		return
	}

	// Construct the http service request
	reqBody, err = json.Marshal(map[string]string{
		"text": prms["text"],
	})
	if err != nil {
		// error constructing an http service request body
		log.Printf("ERROR: %v - error constructing an http service request body. See: %v\n",
			utils.FileLine(),
			err)
		respMap["msg"] = "an error occurred, please try again"
		cMap["havepred"] = false
		respMap["content"] = cMap
		c.JSON(http.StatusInternalServerError, respMap)
		return
	}

	resp, err = http.Post(URLPyMLSvc, "application/json", bytes.NewReader(reqBody))
	if err != nil {
		// error calling the http service
		log.Printf("ERROR: %v - error calling the http service. See: %v\n",
			utils.FileLine(),
			err)
		respMap["msg"] = "an error occurred, please try again"
		cMap["havepred"] = false
		respMap["content"] = cMap
		c.JSON(http.StatusInternalServerError, respMap)
		return
	}
	defer resp.Body.Close()

	rspBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		// error reading the data from the http service response
		log.Printf("ERROR: %v - error reading the data from the http service response. See: %v\n",
			utils.FileLine(),
			err)
		respMap["msg"] = "an error occurred, please try again"
		cMap["havepred"] = false
		respMap["content"] = cMap
		c.JSON(http.StatusInternalServerError, respMap)
		return
	}

	err = utils.FromJSONBytes(rspBody, &respSvc)
	if err != nil {
		// error parsing the data from the http service response
		log.Printf("ERROR: %v - error parsing the data from the http service response. See: %v\n",
			utils.FileLine(),
			err)
		respMap["msg"] = "an error occurred, please try again"
		cMap["havepred"] = false
		respMap["content"] = cMap
		c.JSON(http.StatusInternalServerError, respMap)
		return
	}

	// Return to the caller
	respMap = respSvc

	c.JSON(http.StatusOK, respMap)
}

// genRandMsg generates a random spam or ham text string for use on the frontend
func genRandMsg(c *gin.Context) {
	var idx int
	var class, rtStr string
	var rspMap = make(map[string]string)

	// Determine what type of message is being fetched
	class = c.Param("class")

	// Randomly find a spam message from our original dataset
	if class == "spam" {
		idx = rand.Intn(len(msgStore["spam"]))
		rtStr = msgStore["spam"][idx].Text
	}

	// Randomly find a spam message from our original dataset
	if class == "ham" {
		idx = rand.Intn(len(msgStore["ham"]))
		rtStr = msgStore["ham"][idx].Text
	}

	rspMap["content"] = rtStr
	c.JSON(http.StatusOK, rspMap)
}

func main() {
	stdModel, msgStore, err = fitSpamModel2Data("data/standardSpamData.csv")
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
	r.Static("/vendor", "./www/vendor")
	r.Static("/fonts", "./www/fonts")
	r.Static("/scripts", "./www/scripts")
	r.Static("/views", "./www/views")
	r.Static("/styles", "./www/styles")
	r.Static("/images", "./www/images")
	r.StaticFile("/", "./www/index.html")
	r.StaticFile("/index.html", "./www/index.html")
	r.StaticFile("/index.htm", "./www/index.html")
	r.StaticFile("/favicon.ico", "./www/favicon.ico")
	r.POST("/getModelPred", hdlrPred)
	r.POST("/getPyModelPred", hdlrPyPred)
	r.GET("/getModelInfo", getModelInfo)
	r.GET("/getRandMsg/:class", genRandMsg)

	// Start the web server
	port := "localhost:8090"
	log.Printf("INFO: %v - running web server on port: %v\n",
		utils.FileLine(),
		port)
	r.Run("localhost:8090")
}
