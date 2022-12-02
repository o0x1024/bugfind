package core

import (
	"github.com/tealeg/xlsx"
	"log"
	"strings"
)

func OutputExcel(filename string) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	data := make(map[string]string)
	result, err := LinesInFile(filename)
	if err != nil {
		(err.Error())
		return
	}
	for _, line := range result {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		splits := strings.SplitN(line, " => ", 2)
		if len(splits) != 2 {
			continue
		}
		domain := splits[0]
		ips := strings.Join(strings.Split(splits[1], " => "), ",")
		data[domain] = ips
	}

	log.Println("生成excel..\n")
	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		(err.Error())
		return
	}
	for domain, ips := range data {
		row = sheet.AddRow()
		cell = row.AddCell()
		cell.Value = domain
		cell = row.AddCell()
		cell.Value = ips
	}

	err = file.Save(filename + ".xlsx")
	if err != nil {
		(err.Error())
		return
	}
	log.Println("Excel build success:%s", filename)
}
