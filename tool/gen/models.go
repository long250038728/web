package gen

import (
	"errors"
	"gorm.io/gorm"
	"strings"
	"text/template"
)

type field struct {
	TableSchema string `gorm:"table_schema"`
	TableName   string `gorm:"table_name"`

	Name    string `gorm:"name"`
	Comment string `gorm:"comment"`
	Type    string `gorm:"type"`

	Key string `gorm:"key"`
	Tag string
}

type table struct {
	TableSchema  string   `gorm:"table_schema"`
	TableName    string   `gorm:"table_name"`
	TableComment string   `gorm:"table_comment"`
	Fields       []*field `gorm:"-" json:"-" yaml:"-"`
}

type tableModels struct {
	Tables []*table
}

type Models struct {
	db *gorm.DB
}

func NewModelsGen(db *gorm.DB) *Models {
	return &Models{
		db: db,
	}
}

func (g *Models) Gen(schema string, tables []string) ([]byte, error) {
	if len(tables) == 0 {
		return nil, errors.New("tables num is error")
	}

	list, err := g.dbSearch(schema, tables)
	if err != nil {
		return nil, err
	}

	return (&GenImpl{
		Name:     "gen models",
		TmplPath: "./models.tmpl",
		Func: template.FuncMap{
			"tableName": g.tableName,
			"fieldName": g.fieldName,
			"fieldType": g.fieldType,
		},
		Data: &tableModels{
			Tables: list,
		},
		IsFormat: true,
	}).Gen()
}

func (g *Models) dbSearch(schema string, tables []string) ([]*table, error) {
	var tableList []*table
	var fieldList []*field

	if err := g.db.Debug().Raw(`
	SELECT
		TABLE_SCHEMA as table_schema,
		TABLE_NAME as table_name,
		TABLE_COMMENT as table_comment 
	FROM
		information_schema.TABLES 
	WHERE
		TABLE_SCHEMA = ? AND TABLE_NAME IN (?);
	`, schema, tables).Find(&tableList).Error; err != nil {
		return nil, err
	}
	if len(tableList) != len(tables) {
		return nil, errors.New("search tables count is err")
	}

	if err := g.db.Raw(`
	SELECT
		TABLE_SCHEMA as table_schema,
		TABLE_NAME as table_name,
		COLUMN_NAME as name,
		DATA_TYPE as type,
		COLUMN_COMMENT as comment
	FROM
		information_schema.COLUMNS
	WHERE
		TABLE_SCHEMA = ? AND TABLE_NAME IN (?);
	`, schema, tables).Find(&fieldList).Error; err != nil {
		return nil, err
	}

	for _, table := range tableList {
		for _, fieldItem := range fieldList {
			if table.TableName != fieldItem.TableName {
				continue
			}

			if table.Fields == nil {
				table.Fields = make([]*field, 0, 100)
			}
			table.Fields = append(table.Fields, fieldItem)
		}
	}
	return tableList, nil
}

// tableName 转换表名
func (g *Models) tableName(tableName string) string {
	parts := strings.Split(tableName, "_")
	if len(parts) == 0 {
		return "undefined"
	}
	return g.fieldName(strings.Join(parts[1:len(parts)], "_"))
}

// fieldName 转换字段名
func (g *Models) fieldName(snake string) string {
	// 将字符串分割成数组，以下划线为分隔符
	parts := strings.Split(snake, "_")

	// 遍历数组，将每个部分转换为大写
	var pascal strings.Builder
	for _, part := range parts {
		if len(part) > 0 {
			// 将每个部分的第一个字符转换为大写
			pascal.WriteString(strings.ToUpper(string(part[0])))
			// 将剩余的字符（如果有）添加到结果中
			pascal.WriteString(strings.ToLower(part[1:]))
		}
	}
	// 首字母大写
	return strings.ToUpper(pascal.String()[:1]) + pascal.String()[1:]
}

// fieldType 转换类型
func (g *Models) fieldType(fieldType string) string {
	switch fieldType {
	case "int", "tinyint":
		return "int32" //tinyint int
	case "bigint":
		return "int64" //bigint
	case "decimal":
		return "float32" //decimal
	default:
		return "string" //varchar  char date  datetime json text timestamp
	}
}

var _ = `
-- 获取表基本信息
SELECT
	TABLE_SCHEMA as table_schema,
	TABLE_NAME as table_name,
	TABLE_COMMENT as table_comment 
FROM
	information_schema.TABLES 
WHERE
	TABLE_SCHEMA = "?" AND TABLE_NAME IN (?);
	
-- 获取表信息	
SELECT
	TABLE_SCHEMA as table_schema,
	TABLE_NAME as table_name,
	COLUMN_NAME as name,
	DATA_TYPE as type,
	COLUMN_COMMENT as comment,
	COLUMN_KEY as column_key
FROM
	information_schema.COLUMNS 
WHERE
	TABLE_SCHEMA = "?" AND TABLE_NAME IN (?);

`
