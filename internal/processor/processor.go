package processor

import (
	"context"
	"math"
	"math/rand"
	"os"
	"strconv"

	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"github.com/rs/zerolog"

	res "github.com/riverphillips/envoy-control-plane/api/v1alpha/resources"
	"github.com/riverphillips/envoy-control-plane/internal/watcher"
	"github.com/riverphillips/envoy-control-plane/internal/xdscache"
)

type Processor struct {
	xdsCache        xdscache.XDSCache
	Logger          zerolog.Logger
	snapshotVersion int64
	nodeId          string
	cache           cache.SnapshotCache
}

func (p *Processor) ProcessFile(ctx context.Context, file watcher.NotifyMessage) {
	envoyConfig, err := parseYaml(file.FilePath)
	if err != nil {
		p.Logger.Err(err).Msg("Error parsing yaml")
		return
	}

	for _, l := range envoyConfig.Spec.Listeners {
		var lRoutes []string
		for _, lr := range l.Routes {
			lRoutes = append(lRoutes, lr.Name)
		}

		p.xdsCache.AddListener(l.Name, l.Address, lRoutes, l.Port)

		for _, r := range l.Routes {
			p.xdsCache.AddRoute(r.Name, r.Prefix, r.Clusters)
		}
	}

	for _, c := range envoyConfig.Spec.Clusters {
		p.xdsCache.AddCluster(c.Name)

		for _, e := range c.Endpoints {
			p.xdsCache.AddEndpoint(c.Name, e.Address, e.Port)
		}
	}

	resources := map[resource.Type][]types.Resource{
		resource.RouteType:    p.xdsCache.RouteContents(),
		resource.ClusterType:  p.xdsCache.ClusterContents(),
		resource.ListenerType: p.xdsCache.ListenerContents(),
		resource.EndpointType: p.xdsCache.EndpointContents(),
	}

	snapshot, err := cache.NewSnapshot(
		p.newSnapshotVersion(),
		resources,
	)

	if err != nil {
		p.Logger.Err(err).Msg("Error creating snapshot")
	}

	if err = snapshot.Consistent(); err != nil {
		p.Logger.Err(err).Msg("Inconsistent Snapshot")
		return
	}

	p.Logger.Debug().Interface("Snapshot", snapshot).Msg("Will serve snapshot")

	if err = p.cache.SetSnapshot(ctx, p.nodeId, snapshot); err != nil {
		p.Logger.Err(err).Interface("snapshot", snapshot).Msg("Snapshot error")
		os.Exit(1)
	}
}

func (p *Processor) newSnapshotVersion() string {
	if p.snapshotVersion == math.MaxInt64 {
		p.snapshotVersion = 0
	}

	p.snapshotVersion++
	return strconv.FormatInt(p.snapshotVersion, 10)
}

func New(cache cache.SnapshotCache, nodeId string, logger zerolog.Logger) *Processor {
	return &Processor{
		cache:           cache,
		nodeId:          nodeId,
		snapshotVersion: rand.Int63n(1000),
		Logger:          logger,
		xdsCache: xdscache.XDSCache{
			Listeners: make(map[string]res.Listener),
			Clusters:  make(map[string]res.Cluster),
			Routes:    make(map[string]res.Route),
			Endpoints: make(map[string]res.Endpoint),
		},
	}
}
