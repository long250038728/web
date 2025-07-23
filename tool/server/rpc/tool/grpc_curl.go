package tool

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/fullstorydev/grpcurl"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpCurl struct {
	descriptorSource grpcurl.DescriptorSource
	err              error
}

func NewCurl(importPaths []string, fileNames ...string) *GrpCurl {
	descriptorSource, err := grpcurl.DescriptorSourceFromProtoFiles(importPaths, fileNames...)
	return &GrpCurl{
		descriptorSource: descriptorSource,
		err:              err,
	}
}

func (c *GrpCurl) GetServerMethods() (serverMethods []string, err error) {
	if c.err != nil {
		return nil, c.err
	}

	var servers []string
	if servers, err = grpcurl.ListServices(c.descriptorSource); err != nil {
		return nil, err
	}

	for _, svc := range servers {
		var methods []string
		if methods, err = grpcurl.ListMethods(c.descriptorSource, svc); err != nil {
			return nil, err
		}
		serverMethods = append(serverMethods, methods...)
	}
	return serverMethods, nil
}

func (c *GrpCurl) Curl(ctx context.Context, address, method string, headers map[string]any, data map[string]any) (string, error) {
	if c.err != nil {
		return "", c.err
	}

	b, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	h := make([]string, 0, len(headers))
	for key, value := range headers {
		h = append(h, fmt.Sprintf("%s: %v", key, value))
	}

	// 创建链接
	conn, err := grpcurl.BlockingDial(ctx, "tcp", address, nil, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return "", err
	}
	defer func() {
		_ = conn.Close()
	}()

	// 创建request
	options := grpcurl.FormatOptions{
		EmitJSONDefaultFields: false,
		IncludeTextSeparator:  true,
		AllowUnknownFields:    false,
	}

	requestParser, formatter, err := grpcurl.RequestParserAndFormatter(grpcurl.FormatJSON, c.descriptorSource, bytes.NewReader(b), options)
	if err != nil {
		return "", err
	}
	write := bytes.NewBuffer(nil)
	handler := &grpcurl.DefaultEventHandler{
		Out:            write,
		Formatter:      formatter,
		VerbosityLevel: 0,
	}

	// 调用
	if err = grpcurl.InvokeRPC(ctx, c.descriptorSource, conn, method, h, handler, requestParser.Next); err != nil {
		return "", err
	}
	return write.String(), nil
}
