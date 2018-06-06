package gateway

import (
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type GatewayRegisterer interface {
	RegisterGateway(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error
}
