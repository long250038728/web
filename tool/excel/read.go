package excel

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/long250038728/web/tool/sliceconv"
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
func (r *Read) Read(sheet string, headers []Header, data interface{}) error {
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

	// header 类型检查
	for _, header := range headers {
		if header.Type != HeaderTypeString && header.Type != HeaderTypeInt && header.Type != HeaderTypeFloat && header.Type != HeaderTypeImage {
			return errors.New(fmt.Sprintf("header type error: %s, place use string, int, float, image", header.Type))
		}
	}

	// r.matchHeader(rows[0], headers) 根据第一行的数据，判断每一列的对应的head对象
	return r.marshal(sheet, r.matchHeader(rows[0], headers), rows[1:], data)
}

// matchHeader 根据标题匹配key
func (r *Read) matchHeader(cols []string, headers []Header) []Header {
	//headers 转换为 hash =》 {"name":item}
	var headHash = make(map[string]Header)
	for _, head := range headers {
		headHash[head.Name] = head
	}

	//根据表格中每个col值 ，对传入的header进行排序
	var cellHeader = make([]Header, 0, len(cols))
	for _, headerName := range cols {
		h, ok := headHash[headerName]
		if ok {
			cellHeader = append(cellHeader, h)
		} else {
			cellHeader = append(cellHeader, Header{})
		}
	}
	return cellHeader
}

// marshal 转换为对象
func (r *Read) marshal(sheet string, headers []Header, rows [][]string, data interface{}) error {
	var list = make([]map[string]interface{}, 0, len(rows))

	for rowIndex, row := range rows {
		item := make(map[string]interface{})
		for headIndex, header := range headers {

			val := ""
			if len(row)-1 <= headIndex {
				val = row[headIndex]
			}

			switch header.Type {
			case HeaderTypeString:
				// str类型
				item[header.Key] = val
			case HeaderTypeInt:
				// int类型
				item[header.Key], _ = strconv.ParseInt(val, 10, 0)
			case HeaderTypeFloat:
				// float类型
				item[header.Key], _ = strconv.ParseFloat(val, 32)
			case HeaderTypeImage:
				// image图片
				item[header.Key] = []map[string]any{}
				pics, err := r.file.GetPictures(sheet, fmt.Sprintf("%s%d", cellIndexToName[headIndex], rowIndex+1+1))
				if err == nil || len(pics) > 0 {
					item[header.Key] = sliceconv.Change(pics, func(t excelize.Picture) Pic {
						return Pic{
							File:      t.File,
							Extension: t.Extension,
						}
					})
				}
			}
		}

		//for colIndex, col := range row {
		//	switch headers[colIndex].Type {
		//	case "string":
		//		// str类型
		//		item[headers[colIndex].Key] = col
		//	case "int":
		//		// int类型
		//		item[headers[colIndex].Key], _ = strconv.ParseInt(col, 10, 0)
		//	case "float":
		//		// float类型
		//		item[headers[colIndex].Key], _ = strconv.ParseFloat(col, 32)
		//	default:
		//	}
		//}
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
