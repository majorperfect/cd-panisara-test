package interceptor

import (
	"context"
	"log"

	"google.golang.org/grpc"
)

func Unary(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {

	resp, err := handler(ctx, req)

	return resp, err
}

// Stream returns a server interceptor function to authenticate and authorize stream RPC
func Stream() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		log.Println("--> stream interceptor: ", info.FullMethod)

		return handler(srv, stream)
	}
}
