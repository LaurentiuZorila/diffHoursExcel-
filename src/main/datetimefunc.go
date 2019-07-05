package main

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"
)
const (
	RFC3339FullDate = "2006/01/02 09:00:00 UTC"
	layoutISO = "2006-01-02"
	layoutUS  = "02-02-06"
)

// IsDate returns true when the string is a valid date
func IsDate(str string) bool {
	var fileDate []string
	s := strings.Split(str, "/")
	fileDate = append(fileDate, s[2])
	fileDate = append(fileDate, s[1])
	fileDate = append(fileDate, s[0])
	fullDate := strings.Join(fileDate,"/")
	re := regexp.MustCompile("(0?[1-9]|[12][0-9]|3[01])/(0?[1-9]|1[012])/((19|20)\\d\\d)")
	return re.MatchString(fullDate)
}

func isHour(str string) bool {
	if strings.Contains(str,":") && !strings.Contains(str, "Turno") {
		return true
	}
	return false
}

// return highest or lower value from an array with string values
func getValueByType(arr[]string, low, hight bool) string {
	var value string
	sort.Strings(arr)
	if low {
		value = arr[0]
	} else if hight {
		l := len(arr)
		value = arr[l-1]
	}
	return value
}


func transformDate(str string) string {
	var fileDate []string
	var s []string
	if strings.Contains(str, "/") {
		s = strings.Split(str, "/")
	} else if strings.Contains(str, ".") {
		s = strings.Split(str, ".")
	} else if strings.Contains(str, "-") {
		s = strings.Split(str, "-")
	}

	lengh := len(s)
	year:= s[lengh - 1]
	fullYear := "20" + year
	fileDate = append(fileDate, fullYear)
	fileDate = append(fileDate, s[0])
	fileDate = append(fileDate, s[1])
	fullDate := strings.Join(fileDate,"/")
	//t, _ := time.Parse(layoutISO, fullDate)
	//return t.Format(layoutUS)
	return  fullDate
}

func makeDateAndTime (d, t string) [] string {
	var dateString []string
	var timeString []string
	var dateAndTime [][]string
	var dateStr string
	var timeStr string
	var newTime string
	var tS string
	var tS1 string
	var tS2 string

	t = t + ":00"

	dateString = strings.Split(d,"/")
	timeString = strings.Split(t,":")
	if  strings.Index(timeString[0], "0") == 0 && strings.LastIndex(timeString[0],"0") != 1 {
		//tS = strings.ReplaceAll(timeString[0], "0","")
		tS = strings.Replace(timeString[0],"0","",1)
	} else {
		tS = timeString[0]
	}
	if  strings.Index(timeString[1], "0") == 0 && strings.LastIndex(timeString[1],"0") != 1 {
		//tS1 = strings.ReplaceAll(timeString[0], "0","")
		tS1 = strings.Replace(timeString[1],"0","",1)
	} else {
		tS1 = timeString[1]
	}

	tS2 = strings.Replace(timeString[2],"0","",1)

	newTime = tS + ":" + tS1 + ":" + tS2 + ":0"
	timeString = strings.Split(newTime,":")

	dateAndTime = append(dateAndTime,dateString)
	dateAndTime = append(dateAndTime, timeString)
	dateStr = strings.Join(dateAndTime[0],",")
	timeStr = strings.Join(dateAndTime[1],",")
	aDate := dateStr + "," + timeStr
	return strings.Split(aDate,",")

}

//ShortDateFromString parse shot date from string
func ShortDateFromString(ds string) (time.Time, error) {
	t, err := time.Parse(RFC3339FullDate, ds)
	if err != nil {
		return t, err
	}
	return t, nil
}

//compareDates checks is startdate <= enddate
func compareDates(startdate, enddate string) (bool, error) {
	tstart, err := ShortDateFromString(startdate)
	if err != nil {
		return false, fmt.Errorf("cannot parse startdate: %v", err)
	}
	tend, err := ShortDateFromString(enddate)
	if err != nil {
		return false, fmt.Errorf("cannot parse enddate: %v", err)
	}

	if tstart.After(tend) {
		return false, fmt.Errorf("startdate > enddate - please set proper data boundaries")
	}
	return true, err
}
