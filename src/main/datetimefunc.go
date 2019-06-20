package main

import (
	"fmt"
	"regexp"
	"strconv"
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
	if strings.Contains(str,":") {
		return true
	}
	return false
}

func getLowerValue(arr[]string) string {
	var firstVal int64
	for _,v := range arr {
		t := strings.Split(v,":")
		timeNumber := t[0] + t[1]
		if s, err := strconv.ParseInt(timeNumber, 10, 32); err == nil {
			firstVal = s
		}
	}
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
	fullDate := strings.Join(fileDate,"-")
	//t, _ := time.Parse(layoutISO, fullDate)
	//return t.Format(layoutUS)
	return  fullDate
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
