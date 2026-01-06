package utils

import (
	"context"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/proxy" // import for SOCKS5 support
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type TransportPool struct {
	transports []*http.Transport
	mu         sync.RWMutex
}

func NewTransportPool(proxies []string) *TransportPool {
	pool := &TransportPool{
		transports: make([]*http.Transport, 0, len(proxies)),
	}

	for _, proxyAddr := range proxies {
		var transport *http.Transport
		if strings.HasPrefix(proxyAddr, "socks5://") {
			socks5Proxy, err := url.Parse(proxyAddr)
			if err != nil {
				log.Info().Msgf("Invalid SOCKS5 proxy URL: %v", err)
				continue
			}

			dialer, err := proxy.FromURL(socks5Proxy, proxy.Direct)
			if err != nil {
				log.Info().Msgf("Error creating SOCKS5 dialer: %v", err)
				continue
			}

			transport = &http.Transport{
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					return dialer.Dial(network, addr)
				},
			}
		} else {
			httpProxy, err := url.Parse("http://" + proxyAddr)
			if err != nil {
				log.Info().Msgf("Invalid HTTP proxy URL: %v", err)
				continue
			}

			transport = &http.Transport{
				Proxy: http.ProxyURL(httpProxy),
			}
		}
		pool.transports = append(pool.transports, transport)
	}

	if len(pool.transports) == 0 {
		pool.transports = append(pool.transports, &http.Transport{})
		log.Info().Msg("Found 0 proxy, change to single thread mode")
	}

	return pool
}

func (p *TransportPool) GetRandomTransport() *http.Transport {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.transports[rand.Intn(len(p.transports))]
}

var HttpClientPool = sync.Pool{
	New: func() interface{} {
		return &http.Client{
			Timeout: RequestTimeout * time.Second,
		}
	},
}
