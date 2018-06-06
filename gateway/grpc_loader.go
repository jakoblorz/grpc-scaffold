package gateway

import (
	"log"
	"net/http"

	"github.com/jakoblorz/grpc-scaffold/endpoint"

	"golang.org/x/net/context"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

// Controller represents a grpc controller which
// implements the endpoint.Registerer requirements
// as well as the gateway.Registerer requirements
type Controller interface {
	endpoint.Registerer
	Registerer
}

// GRPCLoader loads all registered controllers once
// Listen() is called - it will create a gateway server
// as well as a grpc server
//
// Calls of RegisterController() will be invoked on the
// endpointLoader as well.
type GRPCLoader struct {
	address        string
	opts           []grpc.DialOption
	controllers    []Registerer
	endpointLoader *endpoint.GRPCLoader
}

// NewGRPCLoader creates a new gateway.GRPCLoader which will
// fork the provided endpoint.GRPCLoader
func NewGRPCLoader(address string, endpointLoader *endpoint.GRPCLoader, opts ...grpc.DialOption) GRPCLoader {

	return GRPCLoader{
		address:        address,
		opts:           opts,
		controllers:    make([]Registerer, 0),
		endpointLoader: endpointLoader,
	}
}

// RegisterController appends the controller to the internal
// controller register
//
// The controller will be registered with the endpoint.GRPCLoader
// as well
func (c GRPCLoader) RegisterController(s Controller) {
	c.controllers = append(c.controllers, s)
	c.endpointLoader.RegisterController(s)
}

// Listen creates both a grpc server and a gateway server
// after having loaded all previously registered controllers
func (c GRPCLoader) Listen() error {

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	mux := runtime.NewServeMux()

	for _, controller := range c.controllers {
		err := controller.RegisterGateway(ctx, mux, c.address, c.opts)
		if err != nil {
			return err
		}
	}

	defer cancel()

	go c.endpointLoader.Listen()

	log.Printf("starting HTTP/1.1 REST server on %s", c.address)

	return http.ListenAndServe(c.address, mux)
}
