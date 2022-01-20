package cmdflag

import (
	"bytes"
	"context"
	"errors"
	errorModel "gRPC_measurement_tool/error"
	"log"
	"runtime"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	// interceptor "gRPC_measurement_tool/interceptors"
	hellowold "gRPC_measurement_tool/protos"
)

func GetGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

func GetInitSetting(cmd Command, startAt time.Time) (uint64, []grpc.DialOption, error) {
	pid := GetGID()
	var opts []grpc.DialOption

	if cmd.IsTls {

		// rootCACert := "../cert/server.crt"
		rootCAKey := "C:\\Users\\danal\\go\\src\\gRPC_measurement_tool\\cert\\rootca.key"
		creds, sslErr := credentials.NewClientTLSFromFile(rootCAKey, "")
		log.Printf("start TLS setting %v", pid)
		errorModel.CheckErrorState(sslErr, pid)

		opts = []grpc.DialOption{
			grpc.WithInsecure(),
			grpc.WithTransportCredentials(creds),
		}
	} else {
		// grpc.WithUnaryInterceptor(interceptor.Identity{ID: pid, StartAt: startAt}.UnaryClient),
		opts = []grpc.DialOption{
			grpc.WithInsecure(),
		}
		log.Printf("start Non TLS setting %v", pid)

	}

	return pid, opts, nil
}

func GetInitTimeout(cmd Command) (context.Context, context.CancelFunc) {
	var ctx context.Context
	var cancel context.CancelFunc
	if cmd.Timeout > 0 {
		ctx, cancel = context.WithTimeout(context.TODO(), time.Duration(cmd.Timeout)*time.Millisecond)
	} else {
		ctx = context.TODO()
		cancel = nil
	}

	return ctx, cancel

}

func GetInitCall(cmd Command, conn *grpc.ClientConn, ctx context.Context) (*hellowold.HelloReply, error) {

	call_slice := strings.Split(cmd.Call, ".")
	if len(call_slice) == 3 {
		package_name := strings.Trim(call_slice[0], " ")
		service_name := strings.Trim(call_slice[1], " ")
		method_name := strings.Trim(call_slice[2], " ")

		switch package_name {
		case "helloworld":
			switch service_name {
			case "Greeter":
				client := hellowold.NewGreeterClient(conn)
				switch method_name {
				case "SayHello":
					r, err := client.SayHello(ctx, &hellowold.HelloRequest{Name: "Mimi"})
					return r, err
				default:
					return nil, errors.New("[404] not found method name")
				}
			default:
				return nil, errors.New("[404] not found service name")
			}
		default:
			return nil, errors.New("[404] not found package name")
		}

		// 서버의 rpc 호출

	} else {
		return nil, errors.New("[400] invalid call value")
	}
}
