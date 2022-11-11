package interceptor

import (
	"context"
	"log"

	"google.golang.org/grpc"
)

// Stream returns a client interceptor to authenticate stream RPC
func Stream() grpc.StreamClientInterceptor {
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption,
	) (grpc.ClientStream, error) {
		log.Printf("--> stream interceptor: %s", method)

		return streamer(ctx, desc, cc, method, opts...)
	}
}
