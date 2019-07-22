package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/tealeg/xlsx"
	_"os"
	"strconv"
	"strings"
	"time"
)

const (
	excelFileName string = "Book1New.xlsx"
	colName string = "Nome"
	colDepartment string = "Sottocommessa"
	colClient string = "Cliente"
	status string = "Stato"
	validStatus string = "IN FORZA"
	preparedFile string = "file1.xlsx"
)

type arrHelper struct {
	notValid, posNotValid bool
}

type validTurn struct {
	turn string
}

type fileHelper struct {
	columns int
	goodTrun int
	nameCol int
	departCol int
	firstDateCol int
	clientCol int
	statusCol int
	fieldsSet bool
}

func main() {

	//t1 := "2019/05/01 09:30:00"
	//t2 := "2019/05/01 11:30:00"
	//diff := diffHours(t1,t2).Hours()
	//if diff > 1 {
	//	fmt.Println("great")
	//}
	//fmt.Println(diff)
	writeNewFile()
}

func prepareFile() [][]string {
	// count columns
	colNumber := countFromSheet(true, false, false, false, excelFileName)
	// total cells
	dim := countFromSheet(false, true, false, true, excelFileName)
	// header file
	var header[] string
	// row
	var r int
	// column
	var c int

	// all values from file
	arr := make([][]string, dim)

	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		fmt.Println(err)
	}
	r = 0
	c = 0
	f := &fileHelper{}
	for _, sheet := range xlFile.Sheets {
		for _, row := range sheet.Rows {
			// values from file
			var arrValues[] string
			for _, cel := range row.Cells {
				if c <= colNumber {
					// set header and columns with name department and status col first row of file
					if r == 0 {
						f.setHeaderColumns(cel.String(), c)
						header = append(header, cel.String())
						arr[r] = header
					} else if r > 0 {
						arrValues = append(arrValues, cel.String())
					}
				}
				c = c + 1
			}
			if checkArr(arrValues, f.statusCol) {
				arr[r] = arrValues
			}
			c = 0
			r = r + 1
		}
	}
	return arr
}

// Contains tells whether a contains x.
func Contains(str string) bool {
	var columnsName = []string {"Nome","Sottocommessa", "Cliente"}
	for _, n := range columnsName {
		if str == n {
			return true
		}
	}
	return false
}

func (f *fileHelper) setHeaderColumns (str string, col int) {
	switch str {
	case colName:
		f.nameCol = col
	case colDepartment:
		f.departCol = col
	case colClient:
		f.clientCol = col
	case status:
		f.statusCol = col
		f.firstDateCol = col + 1
		f.fieldsSet = true
	}
}

func writeNewFile(){
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error
	var counter int
	var request string
	var notValidTurns string = "FERI FERE MALA FEST"
	var possibleNotValidTurns string = "ORNO AANG"
	var riposo string = "Riposo"
	var goodTurn string
	var arrLength int

	// create new file
	file = xlsx.NewFile()
	// add new sheet
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		fmt.Printf(err.Error())
	}

	// add row and cell
	row = sheet.AddRow()
	cell = row.AddCell()

	// array with specific columns
	newFileValues := prepareFile()
	// set fileHelper
	f := &fileHelper{}
	countColumns := len(newFileValues[0])

	// Set file helper
	for c, cel := range newFileValues[0] {
		if f.fieldsSet == false {
			if c <= countColumns {
				// set header and columns with: name, department and status col first row of file
				f.setHeaderColumns(string(cel), c)
			}
		}
		break
	}

	counter = 0
	// Start foreach array with all cells values from file
	for roW, columns := range newFileValues {
		// set length off array on first row (always all cells are completed)
		arrLength = len(columns[0])

		// check if exists values in row and continue
		if len(columns) == 0 {
			continue
		}

		// foreach columns to get cel index => values
		for indx, value := range columns  {
			// if counter greater or equal with number off columns add new row
			if counter == arrLength {
				row = sheet.AddRow()
				cell = row.AddCell()
				counter = 0
			}

			// check if row == 0 and set add in new file cells value with no changes
			if roW == 0 {
				cell.Value = value
				cell = row.AddCell()
				counter = counter + 1
				red := color.New(color.FgRed)
				red.Printf("Cel value is =====> %s   - and counter is %d \n", value, counter)
				continue
			}

			// if cell number is less that first date column
			if indx < f.firstDateCol {
				fmt.Println("c < f.firstDateCol")
				cell.Value = value
				cell = row.AddCell()
				counter = counter + 1
				continue
			}

			// if cell number is equal with first date cel add value
			if f.firstDateCol == indx {
				// if cel value contains "Riposo" add cel value
				if strings.Contains(value, riposo) {
					cell.Value = riposo
					cell = row.AddCell()
					counter = counter + 1
					continue
				}

				// if value is empty add default value
				if len(strings.TrimSpace(value)) == 0 {
					cell.Value = "Not Found Turn"
					cell = row.AddCell()
					counter = counter + 1
					continue
				}

				// check if exist NOT VALID TURNS or POSSIBLE NOT FOUND TUNS
				if strings.Contains(value,riposo) && len(strings.TrimSpace(value)) == 0 {
					var turnValues [] string
					validation := arrHelper{notValid:false, posNotValid:false}
					request = ""

					// Split value and check if is a valid turn
					for _, turns := range strings.Split(value, " ") {
						// if cel value contains not valid turns "FERI MALA ...."
						if strings.Contains(notValidTurns, turns) {
							request = turns
							validation.notValid = true
							continue
						}

						// if cel value contains not valid turns "AANG ORMA ...."
						if strings.Contains(possibleNotValidTurns,turns) {
							// Value of next array value
							val := newFileValues[roW][indx + 1]
							// Check if value is great that zero
							if isGreatThat0(val) {
								validation.posNotValid = false
							} else {
								validation.posNotValid = true
								request = turns
							}
						}
					}

					// if value is a valid turn get highest time
					if validation.notValid == false && validation.posNotValid == false {
						for _, turns := range strings.Split(value, " ") {
							// remove all brackets from value
							turns = removeBrackets(turns)
							// if value is hour and not contains Riposo get first and second time
							if isHour(turns) {
								t := strings.Split(turns, "-")
								turnValues = append(turnValues, t[0])
								turnValues = append(turnValues, t[1])
							}
						}

						// add cell value with highest time
						cell.Value = getTurnHour(turnValues,false,true)
						cell = row.AddCell()
						counter = counter + 1

						// set good turn
						var dateAndTime [] string
						times := transformDate(newFileValues[0][indx]) + " " + getTurnHour(turnValues, false,true)
						// date with highest time "2006/01/02 17:00"
						dateAndTime = append(dateAndTime, times)
						goodTurn = strings.Join(dateAndTime, ",")

						continue
					}

					// get lower time and append turn request
					cell.Value = getTurnHour(turnValues, false,true) + " " + request
					cell = row.AddCell()
					counter = counter + 1
				}
				continue
			}

			// if value == "ORE" continue
			if !checkValidColumns(newFileValues[0][indx]) {
				counter += 1
				continue
			}


			/*
				if cel is > first date col
			 */

			// if cel value contains "Riposo" add cel value
			if strings.Contains(value, riposo) {
				cell.Value = riposo
				cell = row.AddCell()
				counter = counter + 1
				continue
			}

			// if value is empty add default value
			if len(strings.TrimSpace(value)) == 0 {
				cell.Value = "Not Found Turn"
				cell = row.AddCell()
				counter = counter + 1
				continue
			}

			var turnValues [] string
			validation := arrHelper{notValid:false, posNotValid:false}
			request = ""

			// Split value and check if is a valid turn
			for _, turns := range strings.Split(value, " ") {
				// if cel value contains not valid turns "FERI MALA ...."
				if strings.Contains(notValidTurns, turns) {
					request = turns
					validation.notValid = true
					continue
				}

				// if cel value contains not valid turns "AANG ORMA ...."
				if strings.Contains(possibleNotValidTurns,turns) {
					// Value of next array value
					val := newFileValues[roW][indx + 1]
					// Check if value is great that zero
					if isGreatThat0(val) {
						validation.posNotValid = false
					} else {
						validation.posNotValid = true
						request = turns
					}
				}
			}

			// if value is a valid turn get highest time
			if validation.notValid == false && validation.posNotValid == false {
				for _, turns := range strings.Split(value, " ") {
					// remove all brackets from value
					turns = removeBrackets(turns)
					// if value is hour and not contains Riposo get first and second time
					if isHour(turns) {
						t := strings.Split(turns, "-")
						turnValues = append(turnValues, t[0])
						turnValues = append(turnValues, t[1])
					}
				}


				var dateAndTime []string
				var dateAndTimeGoodTurn []string

				times := transformDate(newFileValues[0][indx]) + " " + getTurnHour(turnValues, true, false)
				dateAndTime = append(dateAndTime, times)
				p := strings.Join(dateAndTime,",")
				diff := diffHours(goodTurn, p)

				// add cell value with highest time
				cell.Value = diff.String()
				cell = row.AddCell()
				counter = counter + 1

				times1 := transformDate(newFileValues[0][indx]) + " " + getTurnHour(turnValues, false, true)
				dateAndTimeGoodTurn = append(dateAndTimeGoodTurn, times1)
				p1 := strings.Join(dateAndTimeGoodTurn,times1)
				goodTurn = p1

				continue
			}

			for _, turns := range strings.Split(value, " ") {
				turns = removeBrackets(turns)
				if isHour(turns) {
					t := strings.Split(turns, "-")
					turnValues = append(turnValues, t[0])
					turnValues = append(turnValues,t[1])
				}
			}
			cell.Value = getTurnHour(turnValues, true,false) + " - " + request
			cell = row.AddCell()
			counter += 1
			continue
		}
	}

	// Save file with new data
	err = file.Save(preparedFile)
	if err != nil {
		fmt.Printf(err.Error())
	}
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

// Check if column "stato" is "in forza"
func checkArr(arr[]string, col int) bool {
	if len(arr) > 0 {
		if arr[col] == validStatus {
			return true
		}
	}
	return false
}

// return int with total columns, rows or total cells
func countFromSheet(cel, row, totCells, header bool, file string) int {
	var celNumber int
	var rowNumber int
	var totalCells int
	var value int

	xlFile, err := xlsx.OpenFile(file)
	if err != nil {
		fmt.Println(err)
	}
	for _, sheet := range xlFile.Sheets {
		for r, row := range sheet.Rows {
			rowNumber = r + 1
			for c, _ := range row.Cells {
				celNumber = c + 1
			}
		}
	}

	if header {
		totalCells = rowNumber * celNumber
	} else {
		rowNumber = rowNumber - 1
		totalCells = rowNumber * celNumber
	}
	green := color.New(color.FgWhite)
	magenta := color.New(color.FgHiMagenta)
	if cel {
		value = celNumber

		green.Printf("[%s] ",time.Now().Format("2006-01-02 15:04:05"))
		magenta.Printf("-> Found %d columns on this file: %s\n", value, file)
	} else if row {
		value = rowNumber
		green.Printf("[%s] ",time.Now().Format("2006-01-02 15:04:05"))
		magenta.Printf("-> Found %d rows on this file: %s\n", value, file)
	} else if totCells {
		value = totalCells
		green.Printf("[%s] ",time.Now().Format("2006-01-02 15:04:05"))
		magenta.Printf("-> Found %d total cells on this file: %s \n", value, file)
	}
	return value
}

// check if str contains string "ORE"
func checkValidColumns (str string) bool {
	if strings.Contains(strings.ToLower(str), "ore") {
		return false
	}
	return true
}


