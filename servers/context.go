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

type serverContextStringKey struct{}
type serverContextRequestKey struct{}

type rawKV map[string]string

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

func GetRestRequestCtx(kv ...string) context.Context {
	return GetRequestCtx(REQUEST_TYPE_REST, kv...)
}

func GetJRPCRequestCtx(kv ...string) context.Context {
	return GetRequestCtx(REQUEST_TYPE_JRPC, kv...)
}

func GetTCPRequestCtx(kv ...string) context.Context {
	return GetRequestCtx(REQUEST_TYPE_TCP, kv...)
}

func GetGRPCRequestCtx(kv ...string) context.Context {
	return GetRequestCtx(REQUEST_TYPE_GRPC, kv...)
}

func GetRequestCtx(st RequestType, kv ...string) context.Context {
	newKvs := append(kv, SERVER_REQUEST_TYPE, GetServerTypeValue(st))
	return AppendToRequestCtx(context.Background(), newKvs...)
}

func AppendToRequestCtx(ctx context.Context, kv ...string) context.Context {
	if len(kv)%2 == 1 {
		panic(fmt.Sprintf("metadata: AppendToRequestCtx got an odd number of input pairs for metadata: %d", len(kv)))
	}
	kvs, _ := ctx.Value(serverContextStringKey{}).(rawKV)
	newKvs := make(rawKV, len(kvs)+len(kv)/2)
	for i := 0; i < len(kv); i += 2 {
		newKvs[kv[i]] = kv[i+1]
	}
	for k, v := range kvs {
		newKvs[k] = v
	}
	return context.WithValue(ctx, serverContextStringKey{}, newKvs)
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
	if ctx == nil {
		return map[string]string{}
	}
	t, ok := ctx.Value(serverContextStringKey{}).(rawKV)
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
		res := GetRequestValeFromRaw(raw[0], key)
		return res
	}
	r := GetRequestRaw(ctx)
	if val, ok := r[key]; ok {
		return val
	}
	return ""
}

func GetServerRequestType(ctx context.Context, raw ...map[string]string) string {
	return GetRequestValeByKey(SERVER_REQUEST_TYPE,
		ctx, raw...)
}

func GetServerRequestFunc(ctx context.Context, raw ...map[string]string) string {
	return GetRequestValeByKey(SERVER_REQUEST_FUNC, ctx, raw...)
}

func GetServerRequestInfo(ctx context.Context) string {
	t := ctx.Value(serverContextRequestKey{})
	if t != nil {
		if r, ok := t.(proto.Message); ok {
			return r.String()
		}
	}
	return ""
}

func SetServerRequestInfo(ctx context.Context, req interface{}) context.Context {
	return context.WithValue(ctx, serverContextRequestKey{}, req)
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
	st := GetRequestValeByKey(SERVER_REQUEST_TYPE, ctx)
	ctx = AppendToRequestCtx(ctx, SERVER_REQUEST_FUNC, funcName, SERVER_REQUEST_TYPE, st)
	ctx = SetServerRequestInfo(ctx, req)
	switch st {
	//case REQUEST_TYPE_REST: // 在gincore生成requestIP
	case "grpc":
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
