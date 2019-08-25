package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
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

	data := url.Values{}
	data.Set("from", os.Getenv("FROM"))
	data.Set("to", os.Getenv("TO"))
	data.Set("date", os.Getenv("DATE"))
	data.Set("timetable", os.Getenv("TIMETABLE"))
	data.Set("ticketcount", os.Getenv("TICKETCOUNT"))
	data.Set("carriagecategory", os.Getenv("CARRIAGECATEGORY"))
	data.Set("onlyshowdiscount", os.Getenv("ONLYSHOWDISCOUNT"))
	data.Set("collegestudents", os.Getenv("COLLEGESTUDENTS"))
	data.Set("deviceid", os.Getenv("DEVICEID"))
	data.Set("deviceidhash", os.Getenv("DEVICEIDHASH"))
	data.Set("devicecategory", os.Getenv("DEVICECATEGORY"))
	data.Set("appversion", os.Getenv("APPVERSION"))
	data.Set("parameterversion", os.Getenv("PARAMETERVERSION"))

	req, _ := http.NewRequest("POST", os.Getenv("ENDPOINT"), strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("User-Agent", "ETSiPhone/4.9.5")
	res, _ := http.DefaultClient.Do(req)

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
		restTrains := make([]gjson.Result, 3)
		if len(trains) >= 3 {
			copy(restTrains, trains[0:3])
		} else {
			copy(restTrains, trains)
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
	req, _ := http.NewRequest("POST", os.Getenv("WEBHOOK"), bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}
