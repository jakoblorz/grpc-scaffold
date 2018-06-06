package gateway

import (
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// Registerer defines the requirement for a grpc controller
// to be registerable as gateway reachable endpoint
type Registerer interface {
	RegisterGateway(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error
}
