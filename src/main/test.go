package main

import (
	"fmt"
	"github.com/cheggaaa/pb"
	"github.com/fatih/color"
	"github.com/tealeg/xlsx"
	"strings"
	"time"
)

const (

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

type fileHelper struct {
	columns int
	goodTrun int
	nameCol int
	departCol int
	firstDateCol int
	clientCol int
	statusCol int
	fieldsSet bool
	firstDateValue string
}

func main() {

	//t, _ := time.Parse(layoutISO, fullDate)
	//return t.Format(layoutUS)

	filePath, destinationPath, initError := runInit()
	if !initError {
		count := 3000
		fN := strings.Split(filePath,"/")
		fileName :=  fN[len(fN) - 1]
		bar := pb.StartNew(count).Prefix("Searching file... ")
		bar.ShowCounters = false
		bar.Format("[->_]")
		for i := 0; i < count; i++ {
			bar.Increment()
			time.Sleep(time.Millisecond)
		}
		bar.FinishPrint(" -> File: " + fileName + " has been find")
		bar.Finish()

		color.Blue("Starting to prepare new file...")
		time.Sleep(2 * time.Second)

		writeNewFile(filePath, destinationPath)
		return
	}

	//writeNewFile(excelFileName)
}

func prepareFile(excelFileName string) [][]string {
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
func writeNewFile(excelFile, destination string) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error
	var request string
	var goodTurn string
	var countColumns int
	var newCount int

	// create new file
	file = xlsx.NewFile()
	// add new sheet
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		fmt.Printf(err.Error())
	}

	footerFont := xlsx.NewFont(10, "Verdana")
	footerStyle := xlsx.NewStyle()
	footerStyle.Font = *footerFont

	// array with specific columns
	newFileValues := prepareFile(excelFile)
	// set fileHelper
	f := &fileHelper{}
	countColumns = len(newFileValues[0])

	// Set file helper
	for c, cel := range newFileValues[0] {
		if f.fieldsSet == false {
			if c <= countColumns {
				// set header and columns with: name, department and status col first row of file
				f.setHeaderColumns(string(cel), c)
				// set first date value
				if f.firstDateCol > 0 && len(f.firstDateValue) == 0 {
					f.firstDateValue = newFileValues[0][c + 1]
				}
			}
		}
	}

	// Start progress bar
	barCounter := 0
	// Start foreach array with all cells values from file
	for roW, columns := range newFileValues {
		newCount = 0
		// add row and cell
		if len(columns) > 0 {
			row = sheet.AddRow()
		}

		// if exist values in array
		if len(columns) == 0 {
			continue
		}

		for indx, value := range columns {
			request = ""
			// if header ad only values
			if roW == 0 {
				if checkValidColumns(value) {
					cell = row.AddCell()
					cell.SetStyle(footerStyle)
					cell.Value = value
				}
				continue
			}

			// insert values without changes these are first columns form file (name, id, department ...)
			if newCount < f.firstDateCol {
				cell = row.AddCell()
				cell.SetStyle(footerStyle)
				cell.Value = value
			}

			// insert values form with some changes
			if newCount >= f.firstDateCol && checkValidColumns(newFileValues[0][indx]) {
				// add riposo value contains Riposo
				if hasBreak(value) {
					cell = row.AddCell()
					cell.SetStyle(footerStyle)
					cell.Value = "Riposo"
					continue
				}

				// if is empty value
				if isEmptyString(value) {
					cell = row.AddCell()
					cell.SetStyle(footerStyle)
					cell.Value = "Not found turn"
					continue
				}

				// if exist value and string don't contains "Break"
				var turnValues [] string
				validation := celValue{notValid:false, posNotValid:false}

				if !hasBreak(value) && !isEmptyString(value) {
					// check if value is valid to calculate
					for _, turns := range strings.Split(value, " ") {
						if isNotValidTurn (turns) {
							validation.notValid = true
							request = turns
						}
						if isPossibleNotValidTurn(turns, newFileValues[roW][indx + 1]) {
							validation.posNotValid = true
							request = turns
						}
					}

					// Extract only hours from value
					for _, turns := range strings.Split(value, " ") {
						turns = removeBrackets(turns)
						// append all hours to array
						if isHour(turns) {
							t := strings.Split(turns, "-")
							turnValues = append(turnValues, t[0])
							turnValues = append(turnValues, t[1])
						}
					}

					//add cell value with highest time
					cell = row.AddCell()
					cell.SetStyle(footerStyle)

					if newCount == f.firstDateCol {
						cell.Value = getTurnHour(turnValues,false,true)
						// set good turn
						var dateAndTime [] string
						times := transformDate(newFileValues[0][newCount]) + " " + getTurnHour(turnValues, false,true)
						// date with highest time "2006/01/02 17:00"
						dateAndTime = append(dateAndTime, times)
						goodTurn = strings.Join(dateAndTime, ",")
						//rowGoodTurn = roW
						continue
					}

					if newCount > f.firstDateCol {
						if validation.notValid == false && validation.posNotValid == false {
							var dateAndTime []string
							var dateAndTimeGoodTurn []string
							// set good turn
							times := transformDate(newFileValues[0][indx]) + " " + getTurnHour(turnValues, true, false)
							dateAndTime = append(dateAndTime, times)
							p := strings.Join(dateAndTime,",")

							// diff between good turn and current turn
							diff := diffHours(goodTurn, p)

							if newFileValues[roW][2] == "CUNETE Mihaela" {
								fmt.Println(" ===> good turn: ", goodTurn, " p: ", p, "      ===> diff: ", diff)
							}

							// if dif is less that 12 add red cell value
							if diff.Hours() < 12 {
								redFont := xlsx.NewFont(10, "Verdana")
								redFont.Color = "EC1111"
								redStyle := xlsx.NewStyle()
								redStyle.Font = *redFont
								cell.SetStyle(redStyle)
							}

							// if prev cell contians riposo and diff hour is less that 40 add red cell value
							if strings.Contains(newFileValues[0][newCount - 1], "Riposo") && diff.Hours() < 40 {
								redFont := xlsx.NewFont(10, "Verdana")
								redFont.Color = "EC1111"
								redStyle := xlsx.NewStyle()
								redStyle.Font = *redFont
								cell.SetStyle(redStyle)
							}

							// add cell value with highest time
							cell.Value = diff.String()

							times1 := transformDate(newFileValues[0][indx]) + " " + getTurnHour(turnValues, false, true)
							dateAndTimeGoodTurn = append(dateAndTimeGoodTurn, times1)
							p1 := strings.Join(dateAndTimeGoodTurn,times1)
							goodTurn = p1
						} else {
							cell.Value = request
						}
						continue
					}
				}
				continue
			}
			newCount += 1
			barCounter += 1
			continue
		}
	}

	// Start progress bar
	bar := pb.StartNew(barCounter).Prefix("Loading... ")
	bar.ShowCounters = false
	bar.Format("[=>_]")
	for i := 0; i < barCounter; i++ {
		bar.Increment()
		time.Sleep(time.Millisecond)
	}
	bar.Finish()

	b := color.New(color.FgBlue, color.Bold).SprintFunc()
	fmt.Print(" -> File : ", b(preparedFile), " has been created!\n")
	// Progress bar end

	// Save file with new data
	savedFile := destination + preparedFile
	err = file.Save(savedFile)
	if err != nil {
		red := color.New(color.BgRed)
		red.Printf(err.Error())
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

	if cel {
		value = celNumber
	} else if row {
		value = rowNumber
	} else if totCells {
		value = totalCells
	}
	return value
}




