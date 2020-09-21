package servers

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/legenove/random"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type RequestType int

type serverContextKey struct{}

type rawKV map[string]interface{}

const (
	_ RequestType = iota
	REQUEST_TYPE_REST
	REQUEST_TYPE_GRPC
	REQUEST_TYPE_JRPC
	REQUEST_TYPE_TCP
)

const (
	SERVER_REQUEST_TYPE = "Nano-Request-Type"
	SERVER_REQUEST_FUNC = "Nano-Request-Func"
	SERVER_REQUEST_INFO = "Nano-Request-Info"
)

const (
	SERVER_INCOME_REQUEST_ID   = "Nano-Request-ID"
	SERVER_INCOME_SERVER_NAME  = "Nano-Server-Name"
	SERVER_INCOME_SERVER_GROUP = "Nano-Server-Group"
	SERVER_INCOME_CONTEXT_IP   = "Nano-Context-IP"
	SERVER_INCOME_USER_AGENT   = "User-Agent"
)

func GetRestRequestCtx(kv ...interface{}) context.Context {
	return GetRequestCtx(REQUEST_TYPE_REST, kv...)
}

func GetJRPCRequestCtx(kv ...interface{}) context.Context {
	return GetRequestCtx(REQUEST_TYPE_JRPC, kv...)
}

func GetTCPRequestCtx(kv ...interface{}) context.Context {
	return GetRequestCtx(REQUEST_TYPE_TCP, kv...)
}

func GetGRPCRequestCtx(kv ...interface{}) context.Context {
	return GetRequestCtx(REQUEST_TYPE_GRPC, kv...)
}

func GetRequestCtx(st RequestType, kv ...interface{}) context.Context {
	newKvs := append(kv, SERVER_REQUEST_TYPE, st)
	return AppendToRequestCtx(context.Background(), newKvs...)
}

func AppendToRequestCtx(ctx context.Context, kv ...interface{}) context.Context {
	if len(kv)%2 == 1 {
		panic(fmt.Sprintf("metadata: AppendToRequestCtx got an odd number of input pairs for metadata: %d", len(kv)))
	}
	kvs, _ := ctx.Value(serverContextKey{}).(rawKV)
	newKvs := make(rawKV, len(kvs)+len(kv)/2)
	for i := 0; i < len(kv); i += 2 {
		newKvs[kv[i].(string)] = kv[i+1]
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

func GetRequestRaw(ctx context.Context) map[string]interface{} {
	if ctx == nil {
		return map[string]interface{}{}
	}
	t, ok := ctx.Value(serverContextKey{}).(rawKV)
	if !ok || t == nil {
		return map[string]interface{}{}
	}
	return t
}

func GetRequestValeFromRaw(raw map[string]interface{}, key string) interface{} {
	if val, ok := raw[key]; ok {
		return val
	}
	return nil
}

func GetRequestValeByKey(key string, ctx context.Context, defaultValue interface{}, raw ...map[string]interface{}) interface{} {
	if len(raw) > 0 && raw[0] != nil {
		res := GetRequestValeFromRaw(raw[0], key)
		if res == nil {
			return defaultValue
		}
		return res
	}
	r := GetRequestRaw(ctx)
	if val, ok := r[key]; ok {
		return val
	}
	return defaultValue
}

func GetServerRequestType(ctx context.Context, raw ...map[string]interface{}) string {
	return GetServerTypeValue(GetRequestValeByKey(SERVER_REQUEST_TYPE,
		ctx, REQUEST_TYPE_GRPC, raw...).(RequestType))
}

func GetServerRequestFunc(ctx context.Context, raw ...map[string]interface{}) string {
	return GetRequestValeByKey(SERVER_REQUEST_FUNC, ctx, "", raw...).(string)
}

func GetServerRequestInfo(ctx context.Context, raw ...map[string]interface{}) string {
	req := GetRequestValeByKey(SERVER_REQUEST_INFO, ctx, nil, raw...)
	if req == nil {
		return ""
	}
	if r, ok := req.(proto.Message); ok {
		return r.String()
	}
	return ""
}

// Server Info From MD
func GetServerIncomeByKey(key string, ctx context.Context, raw ...metadata.MD) []string {
	if len(raw) > 0 && raw[0] != nil {
		return raw[0].Get(key)
	}
	r, ok := metadata.FromIncomingContext(ctx)
	if ok {
		return r.Get(key)
	}
	return nil
}

func GetServerName(ctx context.Context, raw ...metadata.MD) string {
	r := GetServerIncomeByKey(SERVER_INCOME_SERVER_NAME, ctx, raw...)
	if r != nil && len(r) > 0 {
		return r[0]
	}
	return ""
}

func GetServerGroup(ctx context.Context, raw ...metadata.MD) string {
	r := GetServerIncomeByKey(SERVER_INCOME_SERVER_GROUP, ctx, raw...)
	if r != nil && len(r) > 0 {
		return r[0]
	}
	return ""
}

func GetRequestId(ctx context.Context, raw ...metadata.MD) string {
	r := GetServerIncomeByKey(SERVER_INCOME_REQUEST_ID, ctx, raw...)
	if r != nil && len(r) > 0 {
		return r[0]
	}
	return random.UuidV5()
}

func GetContextIP(ctx context.Context, raw ...metadata.MD) string {
	r := GetServerIncomeByKey(SERVER_INCOME_CONTEXT_IP, ctx, raw...)
	if r != nil && len(r) > 0 {
		return r[0]
	}
	return ""
}

func GetUserAgent(ctx context.Context, raw ...metadata.MD) string {
	r := GetServerIncomeByKey(SERVER_INCOME_USER_AGENT, ctx, raw...)
	if r != nil && len(r) > 0 {
		return r[0]
	}
	return ""
}

// Server Init Context
func InitContext(ctx context.Context, funcName string, req interface{}) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.MD{}
	}
	st := GetRequestValeByKey(SERVER_REQUEST_TYPE, ctx, REQUEST_TYPE_GRPC).(RequestType)
	ctx = AppendToRequestCtx(ctx, SERVER_REQUEST_FUNC, funcName, SERVER_REQUEST_TYPE, st, SERVER_REQUEST_INFO, req)
	switch st {
	//case REQUEST_TYPE_REST: // 在gincore生成requestIP
	case REQUEST_TYPE_GRPC:
		md.Set(SERVER_INCOME_CONTEXT_IP, RequestIp(ctx))
		//case REQUEST_TYPE_JRPC: // todo
		//case REQUEST_TYPE_TCP: // todo
	}
	if r := md.Get(SERVER_INCOME_REQUEST_ID); r == nil || len(r) == 0 || r[0] == "" {
		md.Set(SERVER_INCOME_REQUEST_ID, random.UuidV5())
		ctx = metadata.NewIncomingContext(ctx, md)
	}
	return ctx
}

func RequestIp(ctx context.Context) string {
	p, _ := peer.FromContext(ctx)
	return p.Addr.String()
}
