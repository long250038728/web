# 代码分析工具链

.PHONY: lint_init
lint_init:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.4.0

.PHONY: lint
lint:
	golangci-lint run ./...