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

// check if string have prefix 01 02 03 04
func checkNotValidHours (s string) bool {
	return strings.HasPrefix(s , "00") || strings.HasPrefix(s , "01") || strings.HasPrefix(s , "02") || strings.HasPrefix(s , "03") || strings.HasPrefix(s , "04")
}


//transform hour 24:00 in 00:00 and if is 00:00 make second param true
func transformNotValidHour (s string) (string, bool) {
	if strings.HasPrefix(s,"24") {
		str := strings.Replace(s,"24", "00", 1)
		return str, true
	} else {
		return s, false
	}
}


// return highest or lower value from an array with string values
func getTurnHour(arr[]string, low, high bool) (string, bool) {
	var value string
	var is24 bool
	sort.Strings(arr)
	if low {
		var newArr []string
		for _, v := range arr {
			if !checkNotValidHours(v) {
				newArr = append(newArr, v)
			}
		}
		sort.Strings(newArr)
		value, is24 = transformNotValidHour(newArr[0])
	} else if high {
		l := len(arr)
		value, is24 = transformNotValidHour(arr[l-1])
	}
	return value + ":00", is24
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

func diffHours(start, end string, t24 bool) time.Duration {
	t1, _ := time.Parse(RFC3339FullDate, start)
	t2, _ := time.Parse(RFC3339FullDate, end)
	var diff time.Duration
	if t24 {
		diff = t2.Sub(t1.Add(time.Hour * 24))
	} else {
		diff = t2.Sub(t1)
	}
	return diff
}

