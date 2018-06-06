package endpoint

import "google.golang.org/grpc"

// Registerer defines the requirement for a grpc controller
// to be registerable as endpoint
type Registerer interface {
	RegisterEndpoint(*grpc.Server)
}
