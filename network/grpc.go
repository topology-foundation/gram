package network

import (
	"context"
	"fmt"
	"net"

	"github.com/topology-gg/gram/config"
	helloPb "github.com/topology-gg/gram/gen/gram/base"
	"google.golang.org/grpc"
)

// TODO: move server implementation to module specific folders
type helloServer struct {
	helloPb.UnimplementedServiceServer
}

func (s *helloServer) SayHello(ctx context.Context, in *helloPb.HelloRequest) (*helloPb.HelloResponse, error) {
	fmt.Printf("Received greet from %s", in.GetName())
	return &helloPb.HelloResponse{Name: "Hello " + in.GetName()}, nil
}

// GRPC represents a struct for GRPC server
type GRPC struct {
	ctx    context.Context
	config *config.GrpcConfig
}

// NewGRPC creates a new grpc server struct
func NewGRPC(ctx context.Context, config *config.GrpcConfig) *GRPC {
	return &GRPC{
		ctx:    ctx,
		config: config,
	}
}

// Start starts a GRPC server
func (g *GRPC) Start() {
	// TODO: read port from config
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", g.config.Port))
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()

	helloPb.RegisterServiceServer(server, &helloServer{})
	fmt.Printf("GRPC Server is listening on port: %v\n", listener.Addr())

	if err := server.Serve(listener); err != nil {
		panic(err)
	}
}
