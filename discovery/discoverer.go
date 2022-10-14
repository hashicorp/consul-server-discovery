package discovery

import (
	"context"
	"net"
	"time"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/go-hclog"
	netaddrs "github.com/hashicorp/go-netaddrs"
)

type Discoverer interface {
	Discover(ctx context.Context) ([]Addr, error)
}

type NetaddrsDiscoverer struct {
	config Config
	log    hclog.Logger
}

func NewNetaddrsDiscoverer(config Config, log hclog.Logger) *NetaddrsDiscoverer {
	return &NetaddrsDiscoverer{
		config: config,
		log:    log,
	}
}

func (n *NetaddrsDiscoverer) Discover(ctx context.Context) ([]Addr, error) {
	start := time.Now()
	addrs, err := netaddrs.IPAddrs(ctx, n.config.Addresses, n.log)
	if err != nil {
		n.log.Error("discovering server addresses", "err", err)
		return nil, err
	}

	metrics.MeasureSince([]string{"discover_servers_duration"}, start)

	var result []Addr
	for _, addr := range addrs {
		result = append(result, Addr{
			TCPAddr: net.TCPAddr{
				IP:   addr.IP,
				Port: n.config.GRPCPort,
				Zone: addr.Zone,
			},
		})
	}
	return result, nil
}
