package endpoint

import "google.golang.org/grpc"

type EndpointRegisterer interface {
	RegisterEndpoint(*grpc.Server)
}
