package main

import (
	"fmt"
	"github.com/cheggaaa/pb"
	"github.com/fatih/color"
	"github.com/tealeg/xlsx"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	colName string = "Nome"
	colDepartment string = "Sottocommessa"
	colClient string = "Cliente"
	colSede string = "Sede"
	status string = "Stato"
	validStatus string = "IN FORZA"
	preparedFile string = "Errors.xlsx"
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
	sedeCol int
	fieldsSet bool
	firstDateValue string
}

func main() {

	fileModified := lastDateModified(path() + fileSerialKey)
	fileKeyLic := readFile(path() + fileSerialKey)
	fileKeyLic = decode(fileKeyLic)
	fileKeyLicArray := strings.Split(fileKeyLic, "___")

	fmt.Println("fileModified -> ", fileModified)
	fmt.Println("fileKeyLicArray[1] -> ", fileKeyLicArray[1])


	os.Exit(3)
	filePath, destinationPath, initError := runInit()

	if !initError {
		// file name
		fileName := fileNameFromPath(filePath)
		// counter for progress bar
		count := 3000

		// yellow color
		warningColor := color.New(color.FgHiYellow).SprintFunc()
		bar := pb.StartNew(count).Prefix(warningColor("Searching file... "))
		bar.ShowCounters = false
		// bar format
		bar.Format("[->_]")
		// loop
		for i := 0; i < count; i++ {
			bar.Increment()
			time.Sleep(time.Millisecond)
		}
		// bar finis msg print
		bar.FinishPrint(warningColor(" -> File: " + fileName + " has been find"))
		bar.Finish()

		infoMsg("Starting to prepare new file...", true)

		// make new file
		writeNewFile(filePath, destinationPath)
		return
	} else {
		if len(filePath) > 0 {
			// file name
			fileName := fileNameFromPath(filePath)
			//counter for progeress bar
			count := 3000
			// yellow color
			info := color.New(color.FgHiYellow).SprintFunc()
			// red color
			danger := color.New(color.FgHiRed).SprintFunc()
			bar := pb.StartNew(count).Prefix(info("Searching file... "))
			bar.ShowCounters = false
			// progress bar format
			bar.Format("[->_]")
			// loop for progess bar
			for i := 0; i < count; i++ {
				bar.Increment()
				time.Sleep(time.Millisecond)
			}
			// msg to print when progress bar finish
			bar.FinishPrint(danger(" Error: -> ") + "File: " + fileName + " doesn't exists, please check again!")
			bar.Finish()
		}
	}
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
	case colSede:
		f.sedeCol = col
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
	var sheet1 *xlsx.Sheet
	var row *xlsx.Row
	var row1 *xlsx.Row
	var cell *xlsx.Cell
	var cell1 *xlsx.Cell
	var err error
	var request string
	var goodTurn string
	var countColumns int
	var newCount int
	var newGoodTime string

	// create new file
	file = xlsx.NewFile()
	// add new sheet
	sheet, err = file.AddSheet("All turns and errors")
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

	var allErrors [][]string


	// Start foreach array with all cells values from file
	for roW, columns := range newFileValues {
		newCount = 0
		goodTurn = ""
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
				barCounter += 1
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
				// get city, department and client names from file
				rowCityName := strings.ToUpper(strings.TrimSpace(newFileValues[roW][f.sedeCol]))
				rowDepartName := strings.ToUpper(strings.TrimSpace(newFileValues[roW][f.clientCol]))
				rowSottName := strings.ToUpper(strings.TrimSpace(newFileValues[roW][f.departCol]))

				// add riposo value contains Riposo
				if hasBreak(value) {
					cell = row.AddCell()
					cell.SetStyle(footerStyle)
					cell.Value = "Riposo"
					barCounter += 1
					continue
				}

				// if is empty value
				if isEmptyString(value) {
					cell = row.AddCell()
					cell.SetStyle(footerStyle)
					cell.Value = "Not found turn"
					barCounter += 1
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
						turnHour, _ := getTurnHour(turnValues, false,true)
						cell.Value = turnHour
						// set good turn
						var dateAndTime [] string
						times := transformDate(newFileValues[0][newCount]) + " " + turnHour
						// date with highest time "2006/01/02 17:00"
						dateAndTime = append(dateAndTime, times)
						goodTurn = strings.Join(dateAndTime, ",")
						barCounter += 1
						continue
					}

					if newCount > f.firstDateCol {
						if validation.notValid == false && validation.posNotValid == false {
							var dateAndTime []string
							var dateAndTimeGoodTurn []string
							// set good turn
							turnHour, _ := getTurnHour(turnValues, true, false)
							// check if hour from good turn is 24:00
							_, t24 := transformNotValidHour(newGoodTime)

							// make date an hour string
							times := transformDate(newFileValues[0][indx]) + " " + turnHour
							dateAndTime = append(dateAndTime, times)
							p := strings.Join(dateAndTime,",")

							// if good turn not exist
							if len(goodTurn) == 0 {
								turnHour, _ := getTurnHour(turnValues, false,true)
								cell.Value = turnHour
								// set good turn
								var dateAndTime [] string
								times := transformDate(newFileValues[0][indx]) + " " + turnHour
								// date with highest time "2006/01/02 17:00"
								dateAndTime = append(dateAndTime, times)
								goodTurn = strings.Join(dateAndTime, ",")
								newGoodTime = turnHour
								barCounter += 1
								continue
							}

							// diff between good turn and current turn
							diff := diffHours(goodTurn, p, t24)
							var rowErrors []string
							// if diff is less that 12 add red cell value
							if diff.Hours() < 12 {
								redFont := xlsx.NewFont(10, "Verdana")
								redFont.Color = "FF0000"
								redStyle := xlsx.NewStyle()
								redStyle.Font = *redFont
								cell.SetStyle(redStyle)

								errorsTableDepartments := rowCityName + " - " + rowDepartName + " - " + rowSottName
								rowErrors = append(rowErrors, errorsTableDepartments, "1")

								if len(allErrors) > 0 {
									if !inArray2D(errorsTableDepartments, allErrors) {
										allErrors = append(allErrors, rowErrors)
									} else {
										for i, v := range allErrors {
											if v[0] == rowErrors[0] {
												errorCounter, _ := strconv.Atoi(v[1])
												errorCounter += 1
												allErrors[i][1] = strconv.Itoa(errorCounter)
											}
										}
									}
								} else {
									allErrors = append(allErrors, rowErrors)
								}
							}

							// if prev cell contains riposo and diff hour is less that 40 add red cell value
							if strings.Contains(newFileValues[0][newCount - 1], "Riposo") && diff.Hours() < 40 {
								redFont := xlsx.NewFont(10, "Verdana")
								redFont.Color = "FF0000"
								redStyle := xlsx.NewStyle()
								redStyle.Font = *redFont
								cell.SetStyle(redStyle)

								errorsTableDepartments := rowCityName + " - " + rowDepartName + " - " + rowSottName
								rowErrors = append(rowErrors, errorsTableDepartments, "1")

								dangerMsg("Errors -> ", false)
								fmt.Println(rowErrors)

								if len(allErrors) > 0 {
									if !inArray2D(errorsTableDepartments, allErrors) {
										allErrors = append(allErrors, rowErrors)
									} else {
										for i, v := range allErrors {
											if v[0] == rowErrors[0] {
												errorCounter, _ := strconv.Atoi(v[1])
												errorCounter += 1
												allErrors[i][1] = strconv.Itoa(errorCounter)
											}
										}
									}
								} else {
									allErrors = append(allErrors, rowErrors)
								}
							}

							// add cell value with highest time
							cell.Value = diff.String()
							turnHour, _ = getTurnHour(turnValues, false, true)
							times1 := transformDate(newFileValues[0][indx]) + " " + turnHour
							dateAndTimeGoodTurn = append(dateAndTimeGoodTurn, times1)
							p1 := strings.Join(dateAndTimeGoodTurn,times1)
							goodTurn = p1
							newGoodTime = turnHour
						} else {
							cell.Value = request
						}
						barCounter += 1
						continue
					}
				}
				barCounter += 1
				continue
			}
			newCount += 1
			barCounter += 1
			continue
		}
	}

	sheet1, err = file.AddSheet("Errors counter")
	if err != nil {
		fmt.Printf(err.Error())
	}

	footerStyle.Font = *footerFont

	row1 = sheet1.AddRow()
	cell1 = row1.AddCell()
	cell1.Value = "CITY"
	cell1 = row1.AddCell()
	cell1.Value = "DEPART"
	cell1 = row1.AddCell()
	cell1.Value = "SOTT"
	cell1 = row1.AddCell()
	cell1.Value = "ERRORS NUMBER"

	for _, v := range allErrors {
		row1 = sheet1.AddRow()

		first := []string{}
		first = strings.Split(v[0],"-")

		for _, cellValue := range first {
			cell1 = row1.AddCell()
			cell1.Value = strings.TrimSpace(cellValue)
		}

		cell1 = row1.AddCell()
		cell1.Value = v[1]

	}

	// Start progress bar
	y := color.New(color.FgHiYellow).SprintFunc()
	bar := pb.StartNew(barCounter).Prefix(y("Writing new file... "))
	bar.ShowCounters = false
	bar.Format("->_")
	for i := 0; i < barCounter; i++ {
		bar.Increment()
		if barCounter <= 10000 {
			time.Sleep(time.Millisecond)
		} else if barCounter > 10000 && barCounter <= 100000 {
			time.Sleep(50 * time.Microsecond)
		} else if barCounter > 100000 && barCounter <= 500000 {
			time.Sleep(10 * time.Microsecond)
		} else if barCounter > 500000 && barCounter <= 1000000 {
			time.Sleep(5 * time.Microsecond)
		} else if barCounter > 100000 && barCounter <= 2000000 {
			time.Sleep(300 * time.Nanosecond)
		} else {
			time.Sleep(time.Nanosecond)
		}
	}
	bar.Finish()
	b := color.New(color.FgHiMagenta, color.Bold).SprintFunc()
	fmt.Print(" -> File : ", b(preparedFile), " has been created!\n")
	fmt.Print(" -> In sheet: ", b("All turns and errors"), " are all turns from file and in red are errors\n")
	fmt.Print(" -> In sheet: ", b("Errors counter"), " are all founded errors\n")

	mySignature := "| by Zorila Laurentiu |"
	sig := "-----------------------------"

	sigColor := color.New(color.FgHiCyan, color.Bold).SprintFunc()

	for i:=0; i<31; i++ {
		fmt.Print(" ")
	}

	for i := 0; i < 23; i++ {
		fmt.Print(sigColor(string(sig[i])))
		time.Sleep(time.Millisecond)
	}

	fmt.Print("\n")

	for i:=0; i<16; i++ {
		fmt.Print(" ")
	}

	for i:=0; i<len(mySignature); i++  {
		fmt.Print(sigColor(string(mySignature[i])))
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Print("\n")

	for i:=0; i<16; i++ {
		fmt.Print(" ")
	}

	for i := 0; i < 23; i++ {
		fmt.Print(sigColor(string(sig[i])))
		time.Sleep(time.Millisecond)
	}

	fmt.Print("\n")


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





