package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/imroc/req"
	"github.com/joho/godotenv"
	"github.com/tidwall/gjson"
)

func main() {
	if os.Getenv("ENV") == "production" {
		lambda.Start(HandleRequest)
	} else {
		HandleRequest()
	}
}

func HandleRequest() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	if os.Getenv("ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			panic("Can't find .env file")
		}
	}

	headers := req.Header{
		"Content-Type": "application/x-www-form-urlencoded",
		"User-Agent":   "ETSiPhone/4.9.5",
	}
	params := req.QueryParam{
		"from":             os.Getenv("FROM"),
		"to":               os.Getenv("TO"),
		"date":             os.Getenv("DATE"),
		"timetable":        os.Getenv("TIMETABLE"),
		"ticketcount":      os.Getenv("TICKETCOUNT"),
		"carriagecategory": os.Getenv("CARRIAGECATEGORY"),
		"onlyshowdiscount": os.Getenv("ONLYSHOWDISCOUNT"),
		"collegestudents":  os.Getenv("COLLEGESTUDENTS"),
		"deviceid":         os.Getenv("DEVICEID"),
		"deviceidhash":     os.Getenv("DEVICEIDHASH"),
		"devicecategory":   os.Getenv("DEVICECATEGORY"),
		"appversion":       os.Getenv("APPVERSION"),
		"parameterversion": os.Getenv("PARAMETERVERSION"),
	}
	r, err := req.Post(os.Getenv("ENDPOINT"), headers, params)
	if err != nil {
		log.Fatal(err)
	}
	res := r.Response()

	defer res.Body.Close()
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	response := string(bodyBytes)

	code := gjson.Get(response, "resultValue.dataStatus")
	message := gjson.Get(response, "resultValue.dataStatusMessage")
	trains := gjson.Get(response, "resultValue.trains").Array()

	if code.String() == "000" && len(trains) != 0 {
		var restTrains []gjson.Result
		for _, v := range trains {
			restTrains = append(restTrains, v)
		}
		obj := map[string]interface{}{}
		for key, train := range restTrains {
			formattedString := createFormattedString(train)
			obj["value"+strconv.Itoa(key+1)] = formattedString
		}
		output, _ := json.Marshal(obj)
		result := sendNotification(output)
		fmt.Println(result)
	} else {
		panic(message.String())
	}
}

func sendNotification(payload []byte) string {
	headers := req.Header{
		"Content-Type": "application/json",
	}
	r, err := req.Post(os.Getenv("WEBHOOK"), headers, req.BodyJSON(payload))
	if err != nil {
		panic(err)
	}
	res := r.Response()
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	return string(body)
}
