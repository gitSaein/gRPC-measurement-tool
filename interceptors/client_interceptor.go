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
	ID      uint64
	StartAt time.Time `json:"start_at"`
}

type ChannelRef struct {
	State              string
	SucceededCount     int32
	TransientFailCount int32
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
	log.Printf("[client-pid: %v] server-status: '%s'", i.ID, cc.GetState())
	md := metadata.Pairs()

	for {
		is_changed_status := cc.WaitForStateChange(ctx, cc.GetState())
		if is_changed_status {
			currentState := cc.GetState()

			elapsed := time.Since(i.StartAt)
			log.Printf("[client-pid: %v] server-status: '%s', take-time: %s, arrive-time: %v", i.ID, currentState, elapsed, time.Now())

			if currentState == connectivity.Ready {
				break
			}
		}

	}

	ctx = metadata.NewOutgoingContext(ctx, md)
	err := invoker(ctx, method, req, reply, cc, opts...)

	return err
}
