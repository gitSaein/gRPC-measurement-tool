package interceptors

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc/metadata"

	"google.golang.org/grpc"
)

type RequestMedatada struct {
	Pid         string
	Method      string
	Authority   string
	ContentType string
	UserAgent   string
}

func parseReq(ctx context.Context) *RequestMedatada {
	md, _ := metadata.FromIncomingContext(ctx)
	method, _ := grpc.Method(ctx)
	pid := md.Get("pid")[0]
	authority := md.Get(":authority")[0]
	contentType := md.Get("content-type")[0]
	userAgent := md.Get("user-agent")[0]
	rm := &RequestMedatada{pid, method, authority, contentType, userAgent}
	log.Printf("%v", rm)
	return rm
}

func UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	parsedReq := parseReq(ctx)

	m, err := handler(ctx, req)
	if err != nil {
		log.Printf("[client-pid: %s][error] server interceptor handler: %v", parsedReq.Pid, err)
	}

	elapsed := time.Since(start)
	log.Printf("[client-pid: %s][take-time: %s] %v ", parsedReq.Pid, elapsed, parsedReq)
	return m, err
}
