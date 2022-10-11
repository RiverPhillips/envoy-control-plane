package xdscache

import (
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"

	res "github.com/riverphillips/envoy-control-plane/api/v1alpha/resources"
	"github.com/riverphillips/envoy-control-plane/internal/resources"
)

type XDSCache struct {
	Listeners map[string]res.Listener
	Endpoints map[string]res.Endpoint
	Routes    map[string]res.Route
	Clusters  map[string]res.Cluster
}

func (xds *XDSCache) AddListener(name, address string, routeNames []string, port uint32) {
	xds.Listeners[name] = res.Listener{
		Name:       name,
		Address:    address,
		Port:       port,
		RouteNames: routeNames,
	}
}

func (xds *XDSCache) AddRoute(name string, prefix string, clusters []string) {
	xds.Routes[name] = res.Route{
		Name:    name,
		Prefix:  prefix,
		Cluster: clusters[0],
	}
}

func (xds *XDSCache) AddCluster(name string) {
	xds.Clusters[name] = res.Cluster{
		Name: name,
	}
}

func (xds *XDSCache) AddEndpoint(clusterName string, address string, port uint32) {
	cluster := xds.Clusters[clusterName]

	cluster.Endpoints = append(cluster.Endpoints, &res.Endpoint{
		UpstreamHost: address,
		UpstreamPort: port,
	})

	xds.Clusters[clusterName] = cluster
}

func (xds *XDSCache) ListenerContents() []types.Resource {
	var r []types.Resource

	for _, l := range xds.Listeners {
		r = append(r, resources.MakeHttpListener(l.Name, l.Address, l.Port))
	}

	return r
}

func (xds *XDSCache) EndpointContents() []types.Resource {
	var r []types.Resource

	for _, c := range xds.Clusters {
		r = append(r, resources.MakeEndpoint(c.Name, c.Endpoints))
	}

	return r
}

func (xds *XDSCache) ClusterContents() []types.Resource {
	var r []types.Resource

	for _, c := range xds.Clusters {
		r = append(r, resources.MakeCluster(c.Name))
	}

	return r
}

func (xds *XDSCache) RouteContents() []types.Resource {
	var routesArr []res.Route

	for _, r := range xds.Routes {
		routesArr = append(routesArr, r)
	}

	return []types.Resource{resources.MakeRoute(routesArr)}
}
