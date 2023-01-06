package proxies

import (
	"strings"
	"sync"
)

type ProxyList struct {
	Name              string
	List              string
	Proxies           []*Proxy
	CurrentProxyIndex int
	Mutex             *sync.Mutex
}

// NewProxyList Creates and returns a new parsed proxy list
func NewProxyList(name string, list string) ProxyList {
	pl := ProxyList{
		Name:              name,
		List:              list,
		CurrentProxyIndex: 0,
	}

	pl.parse()
	return pl
}

// SelectNextProxy Selects the next proxy in the list that is available to use.
// ProxyLists go in order from top to bottom.
func (pl *ProxyList) SelectNextProxy() *Proxy {
	if len(pl.Proxies) == 0 {
		return nil
	}

	pl.Mutex.Lock()
	defer pl.Mutex.Unlock()
	
	proxy := pl.Proxies[pl.CurrentProxyIndex]
	pl.CurrentProxyIndex++

	// Reset the index since the proxy list has already been exhausted
	if pl.CurrentProxyIndex == len(pl.Proxies) {
		pl.CurrentProxyIndex = 0
	}

	return proxy
}

// Parses a raw proxy list
func (pl *ProxyList) parse() {
	pl.Proxies = []*Proxy{}
	pl.Mutex = &sync.Mutex{}
	pl.CurrentProxyIndex = 0

	for _, rawProxy := range strings.Split(strings.ReplaceAll(pl.List, "\r\n", "\n"), "\n") {
		p := Parse(rawProxy)

		if p.Host == "" {
			continue
		}

		pl.Proxies = append(pl.Proxies, &p)
	}
}
