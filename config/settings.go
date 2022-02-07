package config

import (
	"bytes"
	"context"
	"errors"
	"runtime"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"

	m "gRPC_measurement_tool/measure"
	hellowold "gRPC_measurement_tool/protos"
)

func CheckDialConnection(conn *grpc.ClientConn, ctx context.Context, wid uint64, startAt time.Time, report *m.Report) {
	for {
		is_changed_status := conn.WaitForStateChange(ctx, conn.GetState())
		if is_changed_status {
			currentState := conn.GetState()
			elapsed := time.Since(startAt)

			state := &m.ConnectState{ConnectState: currentState, Duration: elapsed, TimeStamp: time.Now()}
			report.States = append(report.States, state)

			if currentState == connectivity.Shutdown || currentState == connectivity.TransientFailure {
				break
			}

		}

	}
}

func GetID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

func SettingOptions(option m.Option) *m.Setting {
	setting := &m.Setting{}

	if option.IsTls {

		rootCACert := "../cert/server.crt"
		creds, err := credentials.NewClientTLSFromFile(rootCACert, "")
		if err != nil {
			return &m.Setting{Error: err}
		}
		setting.Options = append(setting.Options, grpc.WithTransportCredentials(creds))
		setting.Options = append(setting.Options, grpc.FailOnNonTempDialError(true))
		setting.Options = append(setting.Options, grpc.WithBlock())

	} else {
		// grpc.WithUnaryInterceptor(interceptor.Identity{ID: wid, StartAt: startAt}.UnaryClient),
		setting.Options = append(setting.Options, grpc.WithInsecure())
		setting.Options = append(setting.Options, grpc.FailOnNonTempDialError(true))
		setting.Options = append(setting.Options, grpc.WithBlock())
		// grpc.WithStatsHandler(&statsHandler{results: rChan}),

	}

	if option.Timeout > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(option.Timeout)*time.Millisecond)
		setting.Context = ctx
		setting.CancelFunc = cancel
	} else {
		ctx, cancel := context.WithCancel(context.Background())
		setting.Context = ctx
		setting.CancelFunc = cancel

	}

	return setting
}

func CallMethod(option m.Option, conn *grpc.ClientConn, ctx context.Context) (*hellowold.HelloReply, error) {

	if len(option.Call) == 0 {
		return nil, nil
	}

	call_slice := strings.Split(option.Call, ".")
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
