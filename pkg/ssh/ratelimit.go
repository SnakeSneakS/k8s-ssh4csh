package ssh

import (
	"net"
	"strings"
	"time"

	"github.com/gliderlabs/ssh"
	cache "github.com/go-pkgz/expirable-cache/v2"
	"golang.org/x/time/rate"
)

type rateLimiter struct {
	c                cache.Cache[string, *rate.Limiter]
	rateLimitDefault rate.Limit
	rateLimitBurst   int
}

// RateLimiter for each ip address
type RateLimiter interface {
	ConnCallback() ssh.ConnCallback
	ConnectionFailedCallback() ssh.ConnectionFailedCallback
}

func NewRateLimiter() RateLimiter {
	c := cache.NewCache[string, *rate.Limiter]().WithTTL(time.Minute * 10)

	return &rateLimiter{
		c:                c,
		rateLimitDefault: rate.Every(time.Second),
		rateLimitBurst:   1,
	}
}

func (r *rateLimiter) ConnCallback() ssh.ConnCallback {
	return func(ctx ssh.Context, conn net.Conn) net.Conn {
		lim, ok := r.c.Get(getIP(conn.RemoteAddr()))
		if ok {
			lim.Wait(ctx)
		}
		return conn
	}
}

func (r *rateLimiter) ConnectionFailedCallback() ssh.ConnectionFailedCallback {
	return func(conn net.Conn, err error) {
		if err != nil {
			ip := getIP(conn.RemoteAddr())
			_, ok := r.c.Get(ip)
			if !ok {
				lim := rate.NewLimiter(r.rateLimitDefault, r.rateLimitBurst)
				r.c.Set(ip, lim, 0)
				return
			}
			// lim.SetLimit(lim.Limit() / 2)
		}
	}
}

func getIP(addr net.Addr) string {
	ret, _, _ := strings.Cut(addr.String(), ":")
	return ret
}
