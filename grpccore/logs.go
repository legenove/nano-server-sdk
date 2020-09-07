package grpccore

import (
	"context"
	"google.golang.org/grpc"
)

func LoggerRecoveryHandler(funcName string, handler grpc.UnaryHandler) grpc.UnaryHandler {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		// before
		res, err := handler(ctx, req)
		// after
		return res, err
	}
}
