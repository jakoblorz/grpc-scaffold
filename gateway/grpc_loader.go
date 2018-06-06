package gateway

import (
	"log"
	"net/http"

	"github.com/jakoblorz/grpc-scaffold/endpoint"

	"golang.org/x/net/context"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

type Controller interface {
	endpoint.EndpointRegisterer
	GatewayRegisterer
}

type GRPCLoader struct {
	address        string
	opts           []grpc.DialOption
	controllers    []GatewayRegisterer
	endpointLoader *endpoint.GRPCLoader
}

func NewGRPCLoader(address string, endpointLoader *endpoint.GRPCLoader, opts ...grpc.DialOption) GRPCLoader {

	return GRPCLoader{
		address:        address,
		opts:           opts,
		controllers:    make([]GatewayRegisterer, 0),
		endpointLoader: endpointLoader,
	}
}

func (c GRPCLoader) RegisterController(s Controller) {
	c.controllers = append(c.controllers, s)
	c.endpointLoader.RegisterController(s)
}

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
