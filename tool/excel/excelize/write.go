package excelize

import (
	"encoding/json"
	"fmt"
	"github.com/long250038728/web/tool/excel"
	"github.com/xuri/excelize/v2"
)

var SheetDataName = "SheetData"
var cellIndexToName = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

// Write excel工具
type Write struct {
	file *excelize.File
	path string
}

// NewWrite 创建excel工具
func NewWrite(path string) *Write {
	f := excelize.NewFile()
	return &Write{file: f, path: path}
}

// Create 创建excel
func (w *Write) Create(sheet string, headers []excel.Header, data []interface{}) error {
	//标题
	header := make([]interface{}, 0, len(headers))
	for headerIndex, row := range headers {
		header = append(header, row.Name)
		//如果list 有值，会生成一个新的sheet工作区， 把数据填充在里面 ，然后对该行数据进行数据填充
		if err := w.dataSheet(sheet, headerIndex, row.List); err != nil {
			return err
		}
	}
	if err := w.file.SetSheetRow(sheet, fmt.Sprintf("A%d", 1), &header); err != nil {
		return err
	}

	//内容
	for index, row := range data {
		//转换为json
		jsonBytes, _ := json.Marshal(row)
		var jsonData map[string]interface{}
		_ = json.Unmarshal(jsonBytes, &jsonData)

		//值
		cell := make([]interface{}, 0, len(headers)+1)
		for _, h := range headers {
			cell = append(cell, jsonData[h.Key])
		}
		if err := w.file.SetSheetRow(sheet, fmt.Sprintf("A%d", index+2), &cell); err != nil {
			return err
		}
	}
	return w.file.SaveAs(w.path)
}

// dataSheet 如果list 有值，会生成一个新的sheet工作区， 把数据填充在里面 ，然后对该行数据进行数据填充
func (w *Write) dataSheet(sheet string, headerIndex int, list []string) error {
	//如果类型是类别
	if len(list) == 0 {
		return nil
	}

	//还没创建过sheet data 则创建一个
	index, _ := w.file.GetSheetIndex(SheetDataName)
	if index == -1 {
		if _, err := w.file.NewSheet(SheetDataName); err != nil {
			return err
		}
		_ = w.file.SetSheetVisible(SheetDataName, false) //隐藏sheet 如果出错则不显示而已
	}

	cellIndexName := cellIndexToName[headerIndex]
	//在新的sheet对应的列位置把这个枚举输出
	if err := w.file.SetSheetCol(SheetDataName, fmt.Sprintf("%s%d", cellIndexName, 1), &list); err != nil {
		return err
	}

	//使用excl的数据校验功能
	dvRange1 := excelize.NewDataValidation(true)
	dvRange1.Sqref = fmt.Sprintf("%s2:%s10000", cellIndexName, cellIndexName)                                           //需要校验的列
	dvRange1.SetSqrefDropList(fmt.Sprintf("%s!$%s$1:$%s$%d", SheetDataName, cellIndexName, cellIndexName, len(list)+1)) //校验中的枚举列表位置
	return w.file.AddDataValidation(sheet, dvRange1)
}

// Close 关闭excel工具
func (w *Write) Close() {
	_ = w.file.Close()
}
