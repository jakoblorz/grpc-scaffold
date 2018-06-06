package endpoint

import "google.golang.org/grpc"

type Controller interface {
	RegisterEndpoint(*grpc.Server)
}
