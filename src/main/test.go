package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/tealeg/xlsx"

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

type celValue struct {
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
	green := color.New(color.FgRed)
	green.Print("Enter text: ")
	var fileName string
	fmt.Scanln(&fileName)

	if len(fileName) == 0 {
		red := color.New(color.FgRed)
		red.Println("Enter your file name....")
		red.Println("Exit!")
	}
	//writeNewFile()
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

// Set index for name department client and status col
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

// Write new file
func writeNewFile(){
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error
	var request string
	var goodTurn string
	var countColumns int

	// create new file
	file = xlsx.NewFile()
	// add new sheet
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		fmt.Printf(err.Error())
	}

	// array with specific columns
	newFileValues := prepareFile()
	// set fileHelper
	f := &fileHelper{}
	countColumns = len(newFileValues[0])

	// Set file helper
	for c, cel := range newFileValues[0] {
		if f.fieldsSet == false {
			if c <= countColumns {
				// set header and columns with: name, department and status col first row of file
				f.setHeaderColumns(string(cel), c)
			}
		}
	}


	// Start foreach array with all cells values from file
	for roW, columns := range newFileValues {
		// add row and cell
		row = sheet.AddRow()
		// if exist values in array
		if len(columns) == 0 {
			continue
		}

		for indx, value := range columns {
			// if header ad only values
			if roW == 0 && checkValidColumns(value) {
				cell = row.AddCell()
				cell.Value = value
				continue
			}
			// insert values without changes these are first columns form file (name, id, department ...)
			if indx < f.firstDateCol {
				cell = row.AddCell()
				cell.Value = value
			}

			// insert values form with some changes
			if indx >= f.firstDateCol && checkValidColumns(newFileValues[0][indx]) {

				// add riposo value contains Riposo
				if hasBreak(value) {
					cell = row.AddCell()
					cell.Value = "Riposo"
					continue
				}

				// if is empty value
				if isEmptyString(value) {
					cell = row.AddCell()
					cell.Value = "Not found turn"
					continue
				}

				// if exist value and string don't contains "Break"
				var turnValues [] string
				if !hasBreak(value) && !isEmptyString(value) {
					validation := celValue{notValid:false, posNotValid:false}
					// check if value is valid to calculate
					for _, turns := range strings.Split(value, " ") {
						if isNotValidTurn (turns) {
							validation.notValid = true
						}
						if isPossibleNotValidTurn(turns, newFileValues[roW][indx + 1]) {
							validation.posNotValid = true
							request = turns
						}
					}

					// if validation passed
					if validation.notValid == false && validation.posNotValid == false {
						for _, turns := range strings.Split(value, " ") {
							turns = removeBrackets(turns)
							// append all hours to array
							if isHour(turns) {
								t := strings.Split(turns, "-")
								turnValues = append(turnValues, t[0])
								turnValues = append(turnValues, t[1])
							}
						}
						// add cell value with highest time
						cell = row.AddCell()

						//
						if indx == f.firstDateCol {
							cell.Value = getTurnHour(turnValues,false,true)
							// set good turn
							var dateAndTime [] string
							times := transformDate(newFileValues[0][indx]) + " " + getTurnHour(turnValues, false,true)
							// date with highest time "2006/01/02 17:00"
							dateAndTime = append(dateAndTime, times)
							goodTurn = strings.Join(dateAndTime, ",")
							continue
						}

						if indx > f.firstDateCol {
							var dateAndTime []string
							var dateAndTimeGoodTurn []string

							times := transformDate(newFileValues[0][indx]) + " " + getTurnHour(turnValues, true, false)
							dateAndTime = append(dateAndTime, times)
							p := strings.Join(dateAndTime,",")
							diff := diffHours(goodTurn, p)

							// add cell value with highest time
							cell.Value = diff.String()

							times1 := transformDate(newFileValues[0][indx]) + " " + getTurnHour(turnValues, false, true)
							dateAndTimeGoodTurn = append(dateAndTimeGoodTurn, times1)
							p1 := strings.Join(dateAndTimeGoodTurn,times1)
							goodTurn = p1
							continue
						}
					}
					continue
				}

				// if don't exist value or string contains "Break" add cell with default value
				if hasBreak(value) || isEmptyString(value) {
					cell = row.AddCell()
					// get lower time and append turn request
					if indx == f.firstDateCol {
						cell.Value = getTurnHour(turnValues, false, true) + " " + request
					}
					cell.Value = getTurnHour(turnValues, true, false) + " " + request
				}
				continue
			}
		}
	}

	// Save file with new data
	err = file.Save(preparedFile)
	if err != nil {
		fmt.Printf(err.Error())
	}
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




