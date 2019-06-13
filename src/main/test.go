package main

import (
	"fmt"
	"github.com/tealeg/xlsx"
	"regexp"
	"strings"
)

const (
	excelFileName string = "Book1.xlsx"
	colName string = "Nome"
	colDepartment string = "Sottocommessa"
	colClient string = "Cliente"
	status string = "Stato"
	validStatus string = "IN FORZA"
	preparedFile string = "file.xlsx"
	notValid bool = false
	posNotValid bool = false
	RFC3339FullDate = "2006-01-02"
)

type arrHelper struct {
	r ,c int
	allArr [][] string
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

	var notValidTurns string = "FERI FERE MALA"
	//var possibleNotValidTurns string = "ORNO AANG"
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
					if f.firstDateCol >=c {
						if len(strings.Split(cell.Value, " ")) > 1 {
							for _, turns := range strings.Split(cell.Value, " ") {
								if strings.ContainsAny(notValidTurns, turns) {
									// Search in another cel
								} else {
									// Insert turn
								}
							}
						} else if len(strings.Split(cell.Value, " ")) == 2 {
							if strings.ContainsAny(riposo, strings.Split(cell.Value, " ")[1]) {
								// Insert in cel riposo
							} else {
								// Insert in cel turn
							}
						}
					}
					cell.Value = value
					cell = row.AddCell()
					counter = counter + 1
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
func arrValues() [][]string {
	r := countFromSheet(false, true, false, false, excelFileName)
	c := countFromSheet(true, false, false, false, excelFileName)
	var lengs = arrHelper{r: r, c: c}

	arr := make([][]string, r)
	for i := 0; i < lengs.r; i++ {
		arr[i] = make([]string, lengs.c)
	}
	lengs.allArr = arr

	return lengs.allArr
}


