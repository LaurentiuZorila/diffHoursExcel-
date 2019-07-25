package main

import (
	"fmt"
	"strconv"
	"strings"
)

// check if is a empty string
func isEmptyString(s string) bool {
	if len(strings.TrimSpace(s)) == 0 || strings.Contains(s, "Nessun") || strings.Contains(s, "nessun") {
		return true
	}
	return false
}

// check for word "Riposo"
func hasBreak(s string) bool {
	var riposo string = "Riposo"
	if strings.Contains(s,riposo) {
		return true
	}
	return false
}

// check if string contains not valid turns
func isNotValidTurn (s string) bool {
	var notValidTurns string = "FERI FERE MALA FEST"
	s = strings.ToUpper(s)
	if strings.Contains(notValidTurns, s) {
		return true
	}
	return false
}

// check if string contains possible not valid turns
func isPossibleNotValidTurn (s string, d string) bool {
	var possibleNotValidTurns string = "ORNO AANG"
	s = strings.ToUpper(s)
	if strings.Contains(possibleNotValidTurns,s) {
		if isGreatThat0(d) {
			return false
		}
		return true
	}
	return false
}

// check if str contains string "ORE"
func checkValidColumns (str string) bool {
	if strings.Contains(strings.ToLower(str), "ore") {
		return false
	}
	return true
}

// Check if column "stato" is "in forza"
func checkArr(arr[]string, col int) bool {
	if len(arr) > 0 {
		if arr[col] == validStatus {
			return true
		}
	}
	return false
}

// remove all brackets form string
func removeBrackets(str string) string {
	if strings.Count(str, "[") > 0 || strings.Count(str,"]") > 0 {
		str = strings.Replace(str, "[", "", 1)
		str = strings.Replace(str,"]","",1)
		if strings.Count(str, "(") > 0 {
			last := strings.Index(str, "(")
			str = str[0:last]
		}
	} else if strings.HasPrefix(str,"(") {
		if strings.HasSuffix(str,")") {
			str = strings.Replace(str,")", "", 1)
		}
		str = strings.Replace(str,"(", "", 1)
		if strings.Count(str, "(") > 0 {
			last := strings.Index(str, "(")
			str = str[0:last]
		}
	}
	return  str
}

// convert string to int and check if is great that 0
func isGreatThat0 (str string) bool {
	if strings.Contains(str,".") {
		if s, err := strconv.ParseFloat(str, 32); err == nil {
			if s > 0 {
				return true
			}
		} else {
			fmt.Printf("string contains '.' Found errors ======> %s\n", err)
		}
	} else {
		if s, err := strconv.Atoi(str); err == nil {
			if s > 0 {
				return true
			}
		} else {
			fmt.Printf("Found errors ======> %s", err)
		}
	}
	return false
}