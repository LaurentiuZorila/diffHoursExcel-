package main

import (
	"fmt"
	"github.com/tealeg/xlsx"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	excelFileName string = "Book1.xlsx"
	colName string = "Nome"
	colDepartment string = "Sottocommessa"
	colClient string = "Cliente"
	status string = "Stato"
	validStatus string = "IN FORZA"
	preparedFile string = "file.xlsx"
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
	//start := time.Date(2019,03,11,10,0,0,0, time.UTC)
	//end := time.Date(2019,03,11,10,10,10,0, time.UTC)
	//fmt.Println(end.Sub(start))
	year, _ := strconv.Atoi(makeDateAndTime("2019/05/01", "09:30")[0])
	m , _ := strconv.Atoi(makeDateAndTime("2019/05/01", "09:30")[1])
	d , _ := strconv.Atoi(makeDateAndTime("2019/05/01", "09:30")[2])
	hh , _ := strconv.Atoi(makeDateAndTime("2019/05/01", "09:30")[3])
	mm , _ := strconv.Atoi(makeDateAndTime("2019/05/01", "09:30")[4])
	ss , _ := strconv.Atoi(makeDateAndTime("2019/05/01", "09:30")[5])
	ns , _ := strconv.Atoi(makeDateAndTime("2019/05/01", "09:30")[6])

	t1 := time.Date(year,m,d,hh,mm,ss,ns,time.UTC)
	t2 := time.Date(2017, time.February, 16, 0, 0, 0, 0, time.UTC)

	fmt.Println(makeDateAndTime("2019/05/01", "09:30")[0])
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
				// set header and columns with name department and status col first row of file
				f.setHeaderColumns(string(cel), c)
			}
		} else {
			break
		}
	}
	// counter to add new row
	counter = 0
	for r,v := range newFileValues {
		if len(v) > 0 {
			for c, value := range v {
				// if counter greater or equal with number off columns add new row
				if counter == len(v) {
					row = sheet.AddRow()
					cell = row.AddCell()
					counter = 0
				}
				if r == 0 {
					cell.Value = value
					cell = row.AddCell()
					counter = counter + 1
				} else {
					if c >= f.firstDateCol {
						// insert first turn col
						if f.firstDateCol == c {
							if strings.Contains(value,riposo) {
								cell.Value = riposo
								cell = row.AddCell()
								counter = counter + 1
							} else {
								// Check if in cel exist not valid turns or possible not valid turns
								var turnValues [] string
								validation := arrHelper{notValid:false, posNotValid:false}
								request = ""
								for _, turns := range strings.Split(value, " ") {
									if strings.Contains(notValidTurns, turns) {
										request = turns
										validation.notValid = true
									} else if strings.Contains(possibleNotValidTurns, turns) {
										val := newFileValues[r][c + 1]
										// check if string is float or int
										if isGreatThat0(val) {
											validation.posNotValid = false
										} else {
											validation.posNotValid = true
											request = turns
										}
									}
								}
								// if turn is valid get lower time
								if validation.notValid == false && validation.posNotValid == false {
									for _, turns := range strings.Split(value, " ") {
										// check if exist [] in string
										turns = checkCharacters(turns)
										if isHour(turns) {
											t := strings.Split(turns, "-")
											turnValues = append(turnValues, t[0])
											turnValues = append(turnValues, t[1])
										}
									}
									var dateAndTime []string
									times := transformDate(newFileValues[0][c]) + " " + getValueByType(turnValues, true, false)
									dateAndTime = append(dateAndTime, times)

									p := strings.Join(dateAndTime,",")
									fmt.Println(p)
									os.Exit(3)
									fmt.Println(transformDate(newFileValues[0][c]) + " " + getValueByType(turnValues, true, false))
									os.Exit(3)
									cell.Value = getValueByType(turnValues, false, true)
									cell = row.AddCell()
									counter = counter + 1
								} else { // get lower time and append turn request
									cell.Value = getValueByType(turnValues, false, true) + " - " + request
									cell = row.AddCell()
									counter = counter + 1
								}
							}
						} else {
							if strings.Contains(value,riposo) {
								cell.Value = riposo
								cell = row.AddCell()
								counter = counter + 1
							} else {
								// Check if in cel exist not valid turns or possible not valid turns
								var turnValues [] string
								validation := arrHelper{notValid: false, posNotValid: false}
								request = ""
								for _, turns := range strings.Split(value, " ") {
									if strings.Contains(notValidTurns, turns) {
										request = turns
										validation.notValid = true
									} else if strings.Contains(possibleNotValidTurns, turns) {
										val := newFileValues[r][c+1]
										// check if string is float or int
										if isGreatThat0(val) {
											validation.posNotValid = false
										} else {
											validation.posNotValid = true
											request = turns
										}
									}
								}
								// if turn is valid get lower time
								if validation.notValid == false && validation.posNotValid == false {
									for _, turns := range strings.Split(value, " ") {
										// check if exist [] in string
										turns = checkCharacters(turns)
										if isHour(turns) {
											t := strings.Split(turns, "-")
											turnValues = append(turnValues, t[0])
											turnValues = append(turnValues, t[1])
										}
									}
									cell.Value = getValueByType(turnValues, true, false)
									cell = row.AddCell()
									counter = counter + 1

								} else { // get lower time and append turn request
									cell.Value = getValueByType(turnValues, true, false) + " - " + request
									cell = row.AddCell()
									counter = counter + 1
								}
							}
						}
					}
				}
				//fmt.Printf("Counter ===> %d , arrLen ==> %d \n",counter, len(v))
			}
		}
	}
	err = file.Save(preparedFile)
	if err != nil {
		fmt.Printf(err.Error())
	}
}

// convert string to int and check if is great that 0
func isGreatThat0 (str string) bool {
	s, err := strconv.Atoi(str)
	if err == nil {
		if s > 0 {
			return true
		} else {
			return false
		}
	}
	return false
}

func checkCharacters(str string) string {
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


	switch true {
	case cel:
		value = celNumber
	case row:
		value = rowNumber
	case totCells:
		value = totalCells
	}
	return value
}

// return array whit 0 values with number columns and number of cells from file
//func arrValues() [][]string {
//	r := countFromSheet(false, true, false, false, excelFileName)
//	c := countFromSheet(true, false, false, false, excelFileName)
//	var lengs = arrHelper{r: r, c: c}
//
//	arr := make([][]string, r)
//	for i := 0; i < lengs.r; i++ {
//		arr[i] = make([]string, lengs.c)
//	}
//	lengs.allArr = arr
//
//	return lengs.allArr
//}


