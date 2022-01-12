package interceptors

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/metadata"
)

type Identity struct {
	ID      string
	StartAt time.Time `json:"start_at"`
}

func (i Identity) UnaryClient(
	ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	md := metadata.Pairs()
	md.Set("client-pid", i.ID)
	log.Printf("[client-pid: %v] init server status '%s'", i.ID, cc.GetState())
	go func() {
		for {
			if cc.GetState() == connectivity.Ready {
				elapsed := time.Since(i.StartAt)
				log.Printf("[client-pid: %v] changed server status '%s' take-time: %s / at.%v", i.ID, cc.GetState(), elapsed, time.Now())
				break
			}
		}
	}()

	ctx = metadata.NewOutgoingContext(ctx, md)
	err := invoker(ctx, method, req, reply, cc, opts...)

	return err
}
