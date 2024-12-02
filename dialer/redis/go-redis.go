package redis

import (
	"context"
	"net"
	"time"

	"github.com/go-puzzles/puzzles/cores/discover"
	"github.com/go-puzzles/puzzles/plog"
	"github.com/redis/go-redis/v9"
)

func DialGoRedisClient(opts *redis.Options) *redis.Client {
	opts.Dialer = consulGoRedisDial
	return redis.NewClient(opts)
}

func consulGoRedisDial(ctx context.Context, network, addr string) (net.Conn, error) {
	var serviceAddr string

	serviceAddr = discover.GetServiceFinder().GetAddress(addr)
	if serviceAddr == "" {
		serviceAddr = addr
	}
	plog.Debugf("Discover redis addr: %v", serviceAddr)

	netDialer := &net.Dialer{
		Timeout:   time.Second * 5,
		KeepAlive: 5 * time.Minute,
	}

	return netDialer.DialContext(ctx, network, serviceAddr)
}
