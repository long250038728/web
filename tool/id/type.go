package id

type Generate interface {
	Generate() int64
	// GenerateId  isReplace 已经有值是否替换
	GenerateId(model any, isReplace bool) error
}
