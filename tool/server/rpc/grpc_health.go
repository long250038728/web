package rpc

import (
	"context"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type GrpcHealth struct {
	grpc_health_v1.UnimplementedHealthServer
}

func (s *GrpcHealth) Check(ctx context.Context, in *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}, nil
}

func (s *GrpcHealth) Watch(in *grpc_health_v1.HealthCheckRequest, srv grpc_health_v1.Health_WatchServer) error {
	return nil
}
