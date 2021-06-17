package commonutils

import (
	"bufio"
	"fmt"
	"github.com/tealeg/xlsx"
	"os"
	"strings"
)

var (
	xlsxWriterFile    *xlsx.File
	xlsxWriterSheet   *xlsx.Sheet
	xlsxWriterHeadRow *xlsx.Row
	xlsxWriterDataRow *xlsx.Row
	xlsxWriterColumns *xlsx.Col
	xlsxWriterCell    *xlsx.Cell
	excelExtension    string = "xlsx"
)

type XlsxCol struct {
	ColIndex  int
	ColName   string
	ColWidth  float64
	HeadStyle *xlsx.Style
	CellStyle *xlsx.Style
}

type XlsxStyle struct {
	Style *xlsx.Style
}

func NewXlsxStyle(align string, color string, fontName string, size int) *XlsxStyle {
	xlsxStyle := xlsx.NewStyle()
	xlsxStyle.Alignment.Horizontal = align
	xlsxStyle.Font.Underline = false
	xlsxStyle.Fill.BgColor = color
	xlsxStyle.Font.Name = fontName
	xlsxStyle.Font.Size = size
	xlsxStyle.ApplyAlignment = true
	xlsxStyle.ApplyFill = true
	xlsxStyle.ApplyFont = true

	return &XlsxStyle{Style: xlsxStyle}
}




func WriteFileLinesBySplitChar(columns []string, outpath string, input_data []interface{}, splitChar string) bool {
	file, err := os.Create(outpath)
	defer file.Close()
	writer := bufio.NewWriter(file)
	tempData := []string{}
	if len(input_data) > 0 {
		for _, datas := range input_data {
			mapDatas := datas.(map[string]string)
			for _, column := range columns {
				tempData = append(tempData, mapDatas[column])
			}
			_, err := writer.WriteString(strings.Join(tempData, splitChar) + "\n")
			tempData = tempData[0:0]
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	writer.Flush()
	if err == nil {
		return true
	} else {
		return false
	}
}

func WriteFileLinesToExcel(column_tags []XlsxCol, fileName string, input_data []interface{}) bool {
	xlsxWriterFile = xlsx.NewFile()
	sheetName := (strings.Split(fileName, "."))[0]
	xlsxWriterSheet, _ = xlsxWriterFile.AddSheet(sheetName)
	if len(column_tags) > 0 {
		xlsxWriterHeadRow = xlsxWriterSheet.AddRow()
		for _, col := range column_tags {
			xlsxWriterCell = xlsxWriterHeadRow.AddCell()
			xlsxWriterCell.SetStyle(col.HeadStyle)
			xlsxWriterCell.Value = col.ColName
			xlsxWriterSheet.SetColWidth(col.ColIndex, col.ColIndex, col.ColWidth)
		}
	}
	if len(input_data) > 0 {
		for _, datas := range input_data {
			xlsxWriterDataRow = xlsxWriterSheet.AddRow()
			mapDatas := datas.(map[string]string)
			for _, col := range column_tags {
				xlsxWriterCell = xlsxWriterDataRow.AddCell()
				xlsxWriterCell.Value = mapDatas[col.ColName]
				xlsxWriterCell.SetStyle(col.CellStyle)
			}
		}
	}
	err := xlsxWriterFile.Save(fileName)
	if err == nil {
		return true
	} else {
		return false
	}
}
