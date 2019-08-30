package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// check if is a empty string
func isEmptyString(s string) bool {
	//return len(strings.TrimSpace(s)) == 0 || strings.Contains(s, "Nessun") || strings.Contains(s, "nessun")
	if len(strings.TrimSpace(s)) == 0 || strings.Contains(s, "Nessun") || strings.Contains(s, "nessun") {
		return true
	}
	return false
}

// check for word "Riposo"
func hasBreak(s string) bool {
	var riposo = "Riposo"
	//return strings.Contains(s, riposo)
	if strings.Contains(s,riposo) {
		return true
	}
	return false
}

// check if string contains not valid turns
func isNotValidTurn (s string) bool {
	var notValidTurns = "FERI FERE MALA FEST"
	s = strings.ToUpper(s)
	//return strings.Contains(notValidTurns,s)
	if strings.Contains(notValidTurns, s) {
		return true
	}
	return false
}

// check if string contains possible not valid turns
func isPossibleNotValidTurn (s string, d string) bool {
	var possibleNotValidTurns = "ORNO AANG"
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
	//return strings.Contains(strings.ToLower(str), "ore")
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


// run init answers
func runInit () (string, string, bool) {
	var fileName string
	var destination string

	//info.Print(" -> Enter file name and path (ex: C:/user/desktop/etc/file.xlsx): ")
	infoMsg(" -> Enter file name and path (ex: C:/user/desktop/etc/file.xlsx): ", true)
	fmt.Scanln(&fileName)
	excelFileName := fileName

	//info.Print(" -> Enter destination path where save new file (ex: C:/user/desktop/etc/): ")
	infoMsg(" -> Enter destination path where save new file (ex: C:/user/desktop/etc/): ", true)
	fmt.Scanln(&destination)
	destinationPath := destination

	if len(strings.Trim(excelFileName, " ")) == 0 || len(strings.Trim(destinationPath, " ")) == 0 {
		successMsg(" -> File name or destination path missing, please complete all steps!", true)
		//success.Println(" -> File name or destination path missing, please complete all steps....")
		return "", "", true
	} else {
		_, err := os.Stat(excelFileName)
		if err == nil {
			//success.Println(" -> File exists: ", excelFileName)
			if strings.Contains(destinationPath, "/") && !strings.HasSuffix(destinationPath, "/") {
				destinationPath = destinationPath + "/"
			}

			if strings.Contains(destinationPath, "\\") && !strings.HasSuffix(destinationPath, "\\") {
				destinationPath = destinationPath + "\\"
			}

			return excelFileName, destinationPath, false
		} else {
			//danger.Printf(" -> File %s doesn't exists!", fileName)
			return excelFileName, destinationPath, true
		}
	}
}


// return excel file name from path inserted
func fileNameFromPath(filePath string) string {
	var nameOfFile string
	if strings.Contains(filePath, "/") {
		fN := strings.Split(filePath,"/")
		nameOfFile =  fN[len(fN) - 1]
	} else if strings.Contains(filePath, "\\") {
		fN := strings.Split(filePath,"\\")
		nameOfFile =  fN[len(fN) - 1]
	}
	return nameOfFile
}