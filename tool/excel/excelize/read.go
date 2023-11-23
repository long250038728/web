package excelize

import (
	"encoding/json"
	"errors"
	"github.com/long250038728/web/tool/excel"
	"github.com/xuri/excelize/v2"
	"strconv"
)

// Read excel工具
type Read struct {
	file *excelize.File
	err  error
}

// NewRead 新增读excel工具
func NewRead(path string) *Read {
	f, err := excelize.OpenFile(path)
	return &Read{file: f, err: err}
}

// Read 读取excel
func (r *Read) Read(sheet string, headers []excel.Header, data interface{}) error {
	if r.err != nil {
		return r.err
	}

	rows, err := r.file.GetRows(sheet)
	if err != nil {
		return err
	}

	if len(rows) < 2 {
		return errors.New("row num error")
	}

	return r.marshal(r.matchHeader(rows[0], headers), rows[1:], data)
}

// matchHeader 根据标题匹配key
func (r *Read) matchHeader(cols []string, headers []excel.Header) []excel.Header {
	//headers 转换为 hash =》 {"name"：item}
	var headHash = make(map[string]excel.Header)
	for _, head := range headers {
		headHash[head.Name] = head
	}

	//根据表格中每个col值 ，对传入的header进行排序
	var cellHeader = make([]excel.Header, 0, len(cols))
	for _, headerName := range cols {
		h, ok := headHash[headerName]
		if ok {
			cellHeader = append(cellHeader, h)
		} else {
			cellHeader = append(cellHeader, excel.Header{})
		}
	}
	return cellHeader
}

// marshal 转换为对象
func (r *Read) marshal(headers []excel.Header, rows [][]string, data interface{}) error {
	var list = make([]map[string]interface{}, 0, len(rows))

	for _, row := range rows {
		item := make(map[string]interface{})
		for colIndex, col := range row {
			switch headers[colIndex].Type {
			case "string", "list":
				// str类型
				item[headers[colIndex].Key] = col
			case "int":
				// int类型
				item[headers[colIndex].Key], _ = strconv.ParseInt(col, 10, 0)
			case "float":
				// float类型
				item[headers[colIndex].Key], _ = strconv.ParseFloat(col, 32)
			default:
			}
		}
		list = append(list, item)
	}

	jsonBytes, err := json.Marshal(list)
	if err != nil {
		return err
	}

	return json.Unmarshal(jsonBytes, &data)
}

// Close excel关闭文件
func (r *Read) Close() {
	if r.file == nil {
		return
	}
	_ = r.file.Close()
}
