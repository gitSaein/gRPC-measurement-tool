package cmdflag

import (
	"bytes"
	"context"
	"errors"
	"log"
	"runtime"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"

	interceptor "gRPC_measurement_tool/interceptors"
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

func GetInitSetting(cmd Command) (uint64, []grpc.DialOption, time.Time) {
	start := time.Now()
	pid := GetGID()
	log.Printf("[client-pid: %v] arrive-time: %v", pid, start)
	var opts []grpc.DialOption

	if cmd.IsTls {
		opts = []grpc.DialOption{
			grpc.WithUnaryInterceptor(interceptor.Identity{ID: pid, StartAt: start}.UnaryClient),
		}
	} else {
		opts = []grpc.DialOption{
			grpc.WithInsecure(),
			grpc.WithUnaryInterceptor(interceptor.Identity{ID: pid, StartAt: start}.UnaryClient),
		}
	}

	return pid, opts, start
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
