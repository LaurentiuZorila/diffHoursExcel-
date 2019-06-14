package main

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)
const (
	RFC3339FullDate = "2006/01/02 09:00:00 UTC"
)

// IsDate returns true when the string is a valid date
func IsDate(str string) bool {
	var fileDate []string
	s := strings.Split(str, "/")
	lengh := len(s)
	year:= s[lengh - 1]
	fullYear := "20" + year
	fileDate = append(fileDate, s[1])
	fileDate = append(fileDate, s[0])
	fileDate = append(fileDate, fullYear)
	fullDate := strings.Join(fileDate,"/")
	re := regexp.MustCompile("(0?[1-9]|[12][0-9]|3[01])/(0?[1-9]|1[012])/((19|20)\\d\\d)")
	return re.MatchString(fullDate)
}


func transoformDate(str string) string {
	var fileDate []string
	s := strings.Split(str, "/")
	lengh := len(s)
	year:= s[lengh - 1]
	fullYear := "20" + year
	fileDate = append(fileDate, fullYear)
	fileDate = append(fileDate, s[1])
	fileDate = append(fileDate, s[0])
	fullDate := strings.Join(fileDate,"/")
	return fullDate
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
