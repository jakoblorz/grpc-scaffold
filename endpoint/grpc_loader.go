package endpoint

import (
	"log"
	"net"

	"google.golang.org/grpc"
)

// GRPCLoader loads all registered controller
// once Listen() is called
type GRPCLoader struct {
	address     string
	opts        []grpc.ServerOption
	controllers []Registerer
}

// NewGRPCLoader will create a new endpoint.GRPCLoader
func NewGRPCLoader(address string, opts ...grpc.ServerOption) GRPCLoader {

	return GRPCLoader{
		address:     address,
		opts:        opts,
		controllers: make([]Registerer, 0),
	}
}

// RegisterController appends the controller to the internal
// controller register
func (c GRPCLoader) RegisterController(s Registerer) {
	c.controllers = append(c.controllers, s)
}

// Listen creates a new grpc server listening for requests
// after having loaded all previously registerer controllers
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
