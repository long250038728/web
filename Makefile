# 代码分析工具链

.PHONY: lint_init
lint_init:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.4.0


# 需要添加.golangci.yml文件配置相关的插件
.PHONY: lint
lint:
	golangci-lint run ./...