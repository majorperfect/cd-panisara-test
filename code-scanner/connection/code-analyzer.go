package connection

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/majorperfect/guardrails-test/code-scanner/interceptor"

	pb "github.com/majorperfect/guardrails-test/code-analyzer/proto/code-analyzer"
	"github.com/majorperfect/guardrails-test/code-scanner/tls"
	"google.golang.org/grpc"
)

// PaymentMicroservice connect to payment microservice
func CodeAnalyzerMicroservice(c context.Context) (*grpc.ClientConn, pb.CodeAnalyzerServiceClient, error) {
	host := os.Getenv("PAYMENT_HOST")
	port := os.Getenv("PAYMENT_PORT")

	log.Println("admin-gateway", c, "Payment HOST: ", host)
	log.Println("admin-gateway", c, "Payment PORT: ", port)

	addr := host + ":" + port
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
