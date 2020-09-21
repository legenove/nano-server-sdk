package grpccore

import (
	"google.golang.org/grpc"
)

type GrpcDecoratorFunc func(funcName string, handler grpc.UnaryHandler) grpc.UnaryHandler
