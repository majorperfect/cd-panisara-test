package connection

import (
	"context"
	"fmt"
	"log"

	"github.com/majorperfect/guardrails-test/code-scanner/config"

	"github.com/majorperfect/guardrails-test/code-scanner/interceptor"

	pb "github.com/majorperfect/guardrails-test/code-analyzer/proto/code-analyzer"
	"github.com/majorperfect/guardrails-test/code-scanner/tls"
	"google.golang.org/grpc"
)

type AnalyzerMicroserviceFunc = func(ctx context.Context) (*grpc.ClientConn, pb.CodeAnalyzerServiceClient, error)

// PaymentMicroservice connect to payment microservice
func CodeAnalyzerMicroservice(cfg *config.Config) AnalyzerMicroserviceFunc {
	return func(c context.Context) (*grpc.ClientConn, pb.CodeAnalyzerServiceClient, error) {

		log.Println("CodeAnalyzer HOST: ", cfg.ServerHost)
		log.Println("CodeAnalyzer PORT: ", cfg.ServerPort)

		addr := cfg.ServerHost + ":" + cfg.ServerPort
		dialAddr := fmt.Sprintf("dns:///%s", addr)

		options := []grpc.DialOption{}

		cred, err := tls.LoadTLSCredentials()
		if err == nil {
			options = append(options, grpc.WithTransportCredentials(cred))
		} else {
			options = append(options, grpc.WithInsecure())
			log.Println("admin-gateway", c, "Failed to load TLS certificate: %v", err.Error())
		}

		if c != nil {
			options = append(options, grpc.WithStreamInterceptor(interceptor.Stream()))
		}

		conn, err := grpc.Dial(dialAddr, options...)

		if err != nil {
			log.Print("scanner", c, "Can not connect", err.Error())
			return nil, nil, err
		}

		client := pb.NewCodeAnalyzerServiceClient(conn)
		return conn, client, err
	}

}
