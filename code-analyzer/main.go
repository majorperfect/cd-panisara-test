package main

import (
	"log"
	"net"

	"github.com/majorperfect/guardrails-test/code-analyzer/interceptor"
	pb "github.com/majorperfect/guardrails-test/code-analyzer/proto/code-analyzer"
	"github.com/majorperfect/guardrails-test/code-analyzer/server"
	"github.com/majorperfect/guardrails-test/code-analyzer/tls"
	"google.golang.org/grpc"
)

func main() {
	config := Config.SetupConfig()
	addr := config.Host + ":" + config.Port
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	var grpcServer *grpc.Server

	cred, err := tls.LoadTLSCertForServer()
	if err != nil {
		log.Printf("Failed to load TLS certificate: %v", err.Error())
		log.Fatal("cannot load TLS credentials: ", err)

	} else {
		grpcServer = grpc.NewServer(
			grpc.Creds(cred),
			grpc.UnaryInterceptor(interceptor.Unary),
			grpc.StreamInterceptor(interceptor.Stream()),
		)
	}

	pb.RegisterCodeAnalyzerServiceServer(grpcServer, server.New())

	// Serve gRPC Server
	log.Print("Serving gRPC on https://", addr)
	go func() {
		log.Fatal(grpcServer.Serve(lis))
	}()

	grpcServer.Serve(lis)

}
