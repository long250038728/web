package id

type Generate interface {
	Generate() int64
	GenerateId(model any, opts ...Opt) error
}

type GenerateConfig struct {
	fieldName string
	isReplace bool
}

type Opt func(config *GenerateConfig)

// FieldName 字段名称
func FieldName(fieldName string) Opt {
	return func(config *GenerateConfig) {
		config.fieldName = fieldName
	}
}

// IsReplace  已经有值是否替换
func IsReplace(isReplace bool) Opt {
	return func(config *GenerateConfig) {
		config.isReplace = isReplace
	}
}
