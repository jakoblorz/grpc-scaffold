# grpc-scaffold
grpc-scaffold aims to bring a opinionated structure to grpc (+grpc-gateway) repositories.

[![GoDoc](https://godoc.org/github.com/jakoblorz/grpc-scaffold?status.svg)](https://godoc.org/github.com/jakoblorz/grpc-scaffold)
[![Build Status](https://travis-ci.com/jakoblorz/grpc-scaffold.svg?branch=master)](https://travis-ci.com/jakoblorz/grpc-scaffold)
[![codecov](https://codecov.io/gh/jakoblorz/grpc-scaffold/branch/master/graph/badge.svg)](https://codecov.io/gh/jakoblorz/grpc-scaffold)

Each service protobuf definition gets its own struct which may implement the `endpoint.Registerer` or even `gateway.Controller`. By doing so, a common structure is achieved for each service.

## Example
Beginning with the protobuf definition, we define a simple service:
```protobuf
syntax = "proto3";

package proto;

import "google/api/annotations.proto";

service HelloService {
    rpc GetHelloMessage(SayHello) return (ReceiveHello) {
        option (google.api.http) = {
            post: "/api/v1/hello"
            body: "*"
        };
    }
}

message SayHello {
    string user_name = 1;
}

message ReceiveHello {
    string message = 1;
}
```

To implement the service in go, a single struct (one struct per service) is implemented, featuring the functions
required for `endpoint.Registerer` and `gateway.Registerer`.
```go
package main

import (
    "fmt"
    "database/sql"

	"path/to/protobuf/compile/directory/proto"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type HelloService struct {
    Database *sql.DB // db will be initialized later
}

func (h HelloService) GetHelloMessage(c context.Context, request *proto.SayHello) (*proto.ReceiveHello, error) {
    return &proto.ReceiveHello{
        Message: fmt.Sprintf("hello %s", request.UserName),
    }, nil
}

// functions required to implement the gateway.Controller interface:

func (h HelloService) RegisterEndpoint(server *grpc.Server) {
    proto.RegisterHelloServer(server, h)
}

func (h HelloService) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, addr string, opts []grpc.DialOption) error {
    return proto.RegisterHelloHandlerFromEndpoint(ctx, mux, addr, opts)
}
```

During bootstrapping, the service struct is registered as a controller.
```go
package main

import (
    "flag"

    "github.com/jakoblorz/grpc-scaffold/gateway"
    "github.com/jakoblorz/grpc-scaffold/endpoint"

	"google.golang.org/grpc"
)

var (
    grpcPort = flag.Int("grpcport", "8080", "port to listen on with the grpc server")
    gatePort = flag.Int("gateport", "8081", "port to listen on with the gateway")
    host     = flag.String("host", "127.0.0.1", "host to bind to")
)

func main() {
    flag.Parse()


    // db init code
    db := init()

    endpointLoader := endpoint.NewGRPCLoader(fmt.Sprintf("%s:%d", *host, *grpcPort), []grpc.ServerOption{})
    loader := gateway.NewGRPCLoader(fmt.Sprintf("%s:%d", *host, *gatePort), &endpointLoader, []grpc.DialOption{grpc.WithInsecure()})

    loader.RegisterController(HelloService{
        Database: db,
    })

    go loader.Listen()
}
```