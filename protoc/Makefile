.PHONY: init
init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest

.PHONY: gen
gen:
	protoc --proto_path=./ \
      --proto_path=/Users/linlong/go/src/ \
      --go_out=. \
      --go_opt=paths=source_relative \
      --go-grpc_out=. \
      --go-grpc_opt=paths=source_relative \
      --openapi_out=fq_schema_naming=true,default_response=false:. \
      ./*/*.proto