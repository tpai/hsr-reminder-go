package main

import (
	"regexp"

	"github.com/tidwall/gjson"
)

type Station struct {
	Key   string
	Value string
}

var stations = []Station{
	Station{Key: "1", Value: "南港"},
	Station{Key: "2", Value: "台北"},
	Station{Key: "3", Value: "板橋"},
	Station{Key: "4", Value: "桃園"},
	Station{Key: "5", Value: "新竹"},
	Station{Key: "6", Value: "苗栗"},
	Station{Key: "7", Value: "台中"},
	Station{Key: "8", Value: "彰化"},
	Station{Key: "9", Value: "雲林"},
	Station{Key: "10", Value: "嘉義"},
	Station{Key: "11", Value: "台南"},
	Station{Key: "12", Value: "左營"},
}

func getStationById(id string) (station string) {
	for _, v := range stations {
		if v.Key == id {
			return v.Value
		}
	}
	return ""
}

func createFormattedString(data gjson.Result) string {
	var dept string = getStationById(data.Get("deptStation").String())
	var dest string = getStationById(data.Get("destStation").String())
	dateReg := regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)
	timeReg := regexp.MustCompile(`\d{2}:\d{2}:\d{2}`)
	var trainNumber string = data.Get("trainNumber").String()
	var deptDateTime string = data.Get("deptDateTime").String()
	var arrivalDateTime string = data.Get("arrivalDateTime").String()
	var date string = dateReg.FindString(deptDateTime)
	var deptTime string = timeReg.FindString(deptDateTime)
	var arrivalTime string = timeReg.FindString(arrivalDateTime)
	return "[" + trainNumber + "," + dept + "-" + dest + "," + date + "," + deptTime + "-" + arrivalTime + "]"
}
