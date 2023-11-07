package register

import (
	"fmt"
	"math/rand"
)

func HttpServerName(serverName string) string {
	return fmt.Sprintf("%v-%v", serverName, "HTTP")
}

func HttpServerId(serverName string) string {
	return fmt.Sprintf("%v-%v-%d", serverName, "HTTP", rand.Uint64()%10000)
}

func GrpcServerName(serverName string) string {
	return fmt.Sprintf("%v-%v", serverName, "GRPC")
}

func GrpcServerId(serverName string) string {
	return fmt.Sprintf("%v-%v-%d", serverName, "GRPC", rand.Uint64()%10000)
}
