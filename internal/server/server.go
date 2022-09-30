package server

import (
	"context"
	"fmt"
	"log"
	"net"

	clusterv3 "github.com/envoyproxy/go-control-plane/envoy/service/cluster/v3"
	discoveryv3 "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	endpointv3 "github.com/envoyproxy/go-control-plane/envoy/service/endpoint/v3"
	listenerv3 "github.com/envoyproxy/go-control-plane/envoy/service/listener/v3"
	routev3 "github.com/envoyproxy/go-control-plane/envoy/service/route/v3"
	runtimev3 "github.com/envoyproxy/go-control-plane/envoy/service/runtime/v3"
	secretv3 "github.com/envoyproxy/go-control-plane/envoy/service/secret/v3"
	serverv3 "github.com/envoyproxy/go-control-plane/pkg/server/v3"
	"google.golang.org/grpc"
)

const (
	grpcMaxConcurrentStreams = 1_000_000
)

func registerServer(grpcServer *grpc.Server, server serverv3.Server) {
	discoveryv3.RegisterAggregatedDiscoveryServiceServer(grpcServer, server)
	endpointv3.RegisterEndpointDiscoveryServiceServer(grpcServer, server)
	clusterv3.RegisterClusterDiscoveryServiceServer(grpcServer, server)
	listenerv3.RegisterListenerDiscoveryServiceServer(grpcServer, server)
	routev3.RegisterRouteDiscoveryServiceServer(grpcServer, server)
	runtimev3.RegisterRuntimeDiscoveryServiceServer(grpcServer, server)
	secretv3.RegisterSecretDiscoveryServiceServer(grpcServer, server)
}

func Serve(ctx context.Context, srv3 serverv3.Server, port uint) {
	var grpcOptions []grpc.ServerOption
	grpcOptions = append(grpcOptions, grpc.MaxConcurrentStreams(grpcMaxConcurrentStreams))
	grpcServer := grpc.NewServer(grpcOptions...)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}

	registerServer(grpcServer, srv3)

	log.Printf("management server listening on %d\n", port)
	if err = grpcServer.Serve(lis); err != nil {
		log.Println(err)
	}
}
