.PHONY: mac_protobuf
mac_protobuf:
	brew install protobuf

.PHONY: init
init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
	go install github.com/favadi/protoc-go-inject-tag@latest


.PHONY: gen
gen: protoc_gen tag_inject


.PHONY: protoc_gen
protoc_gen:
	protoc --proto_path=./ \
      --proto_path=./proto/ \
      --go_out=. \
      --go_opt=paths=source_relative \
      --go-grpc_out=. \
      --go-grpc_opt=paths=source_relative \
      --openapi_out=fq_schema_naming=true,default_response=false:. \
      ./*/*.proto



.PHONY: tag_inject
tag_inject:
	find . -name "*.pb.go" -exec protoc-go-inject-tag  -remove_tag_comment -input={} \;

