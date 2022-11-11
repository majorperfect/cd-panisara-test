package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/majorperfect/guardrails-test/code-analyzer/interceptor"
	pb "github.com/majorperfect/guardrails-test/code-analyzer/proto/code-analyzer"
	"github.com/majorperfect/guardrails-test/code-analyzer/server"
	"github.com/majorperfect/guardrails-test/code-analyzer/tls"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

type Config struct {
	Host   string `mapstructure:"server_host"`
	Port   string `mapstructure:"server_port"`
	DBHost string `mapstructure:"db_host"`
	DBPort string `mapstructure:"db_port"`
}

func main() {
	config := setupConfig()
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

func setupConfig() *Config {
	viper.SetConfigFile("./config/.env")
	_ = viper.ReadInConfig()

	// environment variable over .env file
	viper.AutomaticEnv()

	cfg := &Config{}
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Panic("unable to decode into config struct, %v \r\n", err)
	}

	///////// BINDING TO GO ENV FOR FUTURE USES /////////
	os.Setenv("ENV", viper.GetString("ENV"))

	os.Setenv("DB_HOST", cfg.DBHost)
	os.Setenv("DB_PORT", cfg.DBPort)
	os.Setenv("SERVER_HOST", cfg.Host)
	os.Setenv("SERVER_PORT", cfg.Port)

	// TLS certificate
	os.Setenv("TLS_CERT_FILE", viper.GetString("TLS_CERT_FILE"))
	os.Setenv("TLS_KEY_FILE", viper.GetString("TLS_KEY_FILE"))
	////////////////////////////////////////////////////
	// os.Setenv("JWT_KEY", viper.GetString("JWT_KEY"))
	for _, element := range os.Environ() {
		variable := strings.Split(element, "=")
		fmt.Println(variable[0], "=>", variable[1])
	}

	return cfg
}
