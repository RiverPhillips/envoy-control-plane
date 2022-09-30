package xdscache

import (
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"

	"github.com/riverphillips/envoy-control-plane/internal/resources"
)

type XDSCache struct {
	Listeners map[string]resources.Listener
	Endpoints map[string]resources.Endpoint
	Routes    map[string]resources.Route
	Clusters  map[string]resources.Cluster
}

func (c XDSCache) AddListener(name string, routes []string, address string, port uint32) {

}

func (c XDSCache) AddRoute(name string, prefix string, names []string) {

}

func (c XDSCache) AddCluster(name string) {

}

func (c XDSCache) AddEndpoint(name string, address string, port uint32) {

}

func (c XDSCache) ListenerContents() types.Resource {

}
