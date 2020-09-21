package servers

import (
	"context"
	"fmt"
)

type RequestType int

type serverContextKey struct{}

type rawKV map[string]string

const (
	_ RequestType = iota
	REQUEST_TYPE_REST
	REQUEST_TYPE_GRPC
	REQUEST_TYPE_JRPC
	REQUEST_TYPE_TCP
)

const (
	SERVER_REQUEST_TYPE_KEY = "Nano-Request-Type"
	SERVER_REQUEST_FUNC_KEY = "Nano-Request-Func"
)

func GetRestRequestCtx() context.Context {
	return GetRequestCtx(REQUEST_TYPE_REST)
}

func GetJRPCRequestCtx() context.Context {
	return GetRequestCtx(REQUEST_TYPE_JRPC)
}

func GetTCPRequestCtx() context.Context {
	return GetRequestCtx(REQUEST_TYPE_TCP)
}

func GetGRPCRequestCtx() context.Context {
	return GetRequestCtx(REQUEST_TYPE_GRPC)
}

func GetRequestCtx(st RequestType) context.Context {
	newKvs := rawKV{SERVER_REQUEST_TYPE_KEY: GetServerTypeValue(st)}
	return context.WithValue(context.Background(), serverContextKey{}, newKvs)
}

func AppendToRequestCtx(ctx context.Context, kv ...string) context.Context {
	if len(kv)%2 == 1 {
		panic(fmt.Sprintf("metadata: AppendToRequestCtx got an odd number of input pairs for metadata: %d", len(kv)))
	}
	kvs, _ := ctx.Value(serverContextKey{}).(rawKV)
	newKvs := make(rawKV, len(kvs)+len(kv)/2)
	for i := 0; i < len(kv); i += 2 {
		newKvs[kv[i]] = kv[i+1]
	}
	for k, v := range kvs {
		newKvs[k] = v
	}
	return context.WithValue(ctx, serverContextKey{}, newKvs)
}

func GetServerTypeValue(st RequestType) string {
	switch st {
	case REQUEST_TYPE_REST:
		return "rest"
	case REQUEST_TYPE_GRPC:
		return "grpc"
	case REQUEST_TYPE_JRPC:
		return "jrpc"
	case REQUEST_TYPE_TCP:
		return "tcp"
	}
	return "grpc"
}

func GetRequestRaw(ctx context.Context) map[string]string {
	t, ok := ctx.Value(serverContextKey{}).(rawKV)
	if !ok || t == nil {
		return map[string]string{}
	}
	return t
}

func GetRequestValeFromRaw(raw map[string]string, key string) string {
	if val, ok := raw[key]; ok {
		return val
	}
	return ""
}

func GetRequestValeByKey(key string, ctx context.Context, raw ...map[string]string) string {
	if len(raw) > 0 && raw[0] != nil {
		return GetRequestValeFromRaw(raw[0], key)
	}
	r := GetRequestRaw(ctx)
	if val, ok := r[key]; ok {
		return val
	}
	return ""
}

func GetServerRequestType(ctx context.Context, raw ...map[string]string) string {
	return GetRequestValeByKey(SERVER_REQUEST_TYPE_KEY, ctx, raw...)
}

func GetServerRequestFunc(ctx context.Context, raw ...map[string]string) string {
	return GetRequestValeByKey(SERVER_REQUEST_FUNC_KEY, ctx, raw...)
}

// Server Info From MD

