package main

import (
	"sort"
	"strings"
	"time"
)
const (
	RFC3339FullDate = "2006/01/02 15:04:05"
)

func isHour(str string) bool {
	if strings.Contains(str,":") && !strings.Contains(str, "Turno") {
		return true
	}
	return false
}

// return highest or lower value from an array with string values
func getTurnHour(arr[]string, low, high bool) string {
	var value string
	sort.Strings(arr)
	if low {
		value = arr[0]
	} else if high {
		l := len(arr)
		value = arr[l-1]
	}
	return value + ":00"
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

	length := len(s)
	year:= s[length - 1]
	fullYear := "20" + year
	fileDate = append(fileDate, fullYear)
	fileDate = append(fileDate, s[0])
	fileDate = append(fileDate, s[1])
	fullDate := strings.Join(fileDate,"/")
	//t, _ := time.Parse(layoutISO, fullDate)
	//return t.Format(layoutUS)

	return fullDate
}

func diffHours(start, end string) time.Duration {
	t1, _ := time.Parse(RFC3339FullDate, start)
	t2, _ := time.Parse(RFC3339FullDate, end)
	diff := t2.Sub(t1)
	return  diff
}

