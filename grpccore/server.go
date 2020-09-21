package grpccore

import (
	"sync"

	"google.golang.org/grpc"
)

type RegisterServer func(s *grpc.Server)

var registerMapper = map[string]RegisterServer{}
var rpcServer *grpc.Server
var mu sync.Mutex

func GetServerWithOptions(opt ...grpc.ServerOption) *grpc.Server {
	mu.Lock()
	defer mu.Unlock()
	if rpcServer == nil {
		rpcServer = grpc.NewServer(opt...)
	}
	for _, f := range registerMapper {
		f(rpcServer)
	}
	return rpcServer
}

func RegisterToServer(n string, f RegisterServer) {
	if registerMapper == nil {
		registerMapper = map[string]RegisterServer{}
	}
	registerMapper[n] = f
}
