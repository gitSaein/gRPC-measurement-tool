package rpc

import (
	"context"
	pb "grpc/protos/health"
	proto "grpc/protos/health"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type healthClient struct {
	client pb.HealthClient
	conn   *grpc.ClientConn
}

func NewGrpcHealthClient(conn *grpc.ClientConn) Health {
	client := new(healthClient)
	client.client = pb.NewHealthClient(conn)
	client.conn = conn
	return client
}

func (c *healthClient) Close() error {
	return c.conn.Close()
}

func (c *healthClient) Check(ctx context.Context) (bool, error) {
	var res *proto.HealthCheckResponse
	var err error
	req := new(proto.HealthCheckRequest)

	res, err = c.client.Check(ctx, req)
	if err == nil {
		if res.GetStatus() == proto.HealthCheckResponse_SERVING {
			return true, nil
		}
		return false, nil
	}
	log.Fatalf("[error healthcheck] %v", grpc.Code(err))
	switch grpc.Code(err) {
	case codes.Aborted,
		codes.DataLoss,
		codes.DeadlineExceeded,
		codes.Internal,
		codes.Unavailable:
	default:
		return false, err

	}
	return false, err

}
