package orm

import (
	"errors"
	"fmt"
	"github.com/xwb1989/sqlparser"
	"os"
	"strings"
)

func (g *Gorm) Parser(file string) ([]string, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil, err
	}

	// 将文件内容转换为字符串
	sqlContent := string(content)

	// 使用 strings.Split 分解SQL语句
	statements := strings.Split(sqlContent, ";")

	var sqls = make([]string, 0, len(statements))

	// 解析每个SQL语句
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}

		parsedStmt, err := sqlparser.Parse(stmt)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("sql parser is err %s, stmt:%s", err.Error(), stmt))
		}
		switch parsedStmt.(type) {
		case *sqlparser.DDL:
			sqls = append(sqls, stmt)
		case *sqlparser.Insert, *sqlparser.Update, *sqlparser.Delete:
			sqls = append(sqls, sqlparser.String(parsedStmt))
		default:
			return nil, errors.New("sql parsed stmt type undefine")
		}
	}
	return sqls, nil
}
