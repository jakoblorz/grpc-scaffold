package endpoint

import (
	"log"
	"net"

	"google.golang.org/grpc"
)

type GRPCLoader struct {
	address     string
	opts        []grpc.ServerOption
	controllers []EndpointRegisterer
}

func NewGRPCLoader(address string, opts ...grpc.ServerOption) GRPCLoader {

	return GRPCLoader{
		address:     address,
		opts:        opts,
		controllers: make([]EndpointRegisterer, 0),
	}
}

func (c GRPCLoader) RegisterController(s EndpointRegisterer) {
	c.controllers = append(c.controllers, s)
}

func (c GRPCLoader) Listen() error {

	server := grpc.NewServer(c.opts...)

	for _, controller := range c.controllers {
		controller.RegisterEndpoint(server)
	}

	listen, err := net.Listen("tcp", c.address)
	if err != nil {
		return err
	}

	log.Printf("starting gRPC server on %s", c.address)

	return server.Serve(listen)
}
