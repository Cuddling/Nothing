package proxies

import (
	"fmt"
	"sync"
	"testing"
)

func TestParseProxyList(t *testing.T) {
	raw := "127.0.0.1:8080:username:password\r\n" +
		"192.68.1.1:1234\r\n" +
		"255.255.255.255:0000:user:pass"

	list := NewProxyList("Test", raw)

	if len(list.Proxies) != 3 {
		t.Fatalf("expected 3 proxies. got %v", len(list.Proxies))
	}
}

func TestSelectProxies(t *testing.T) {
	raw := "127.0.0.1:8080:username:password\r\n" +
		"192.68.1.1:1234\r\n" +
		"255.255.255.255:0000:user:pass"

	list := NewProxyList("Test", raw)

	p1 := list.SelectNextProxy()
	p2 := list.SelectNextProxy()
	p3 := list.SelectNextProxy()

	if p1 != list.Proxies[0] {
		t.Fatalf("incorrect proxy at index zero")
	}

	if p2 != list.Proxies[1] {
		t.Fatalf("incorrect proxies at index one")
	}

	if p3 != list.Proxies[2] {
		t.Fatalf("incorrect proxies at index two")
	}

	fmt.Println(p1)
	fmt.Println(p2)
	fmt.Println(p3)
}

func TestSelectProxiesGoRoutines(t *testing.T) {
	raw := ""
	const proxyCount int = 100

	for i := 0; i < proxyCount; i++ {
		raw += fmt.Sprintf("127.0.0.%v:8080\r\n", proxyCount)
	}

	list := NewProxyList("Test", raw)

	if len(list.Proxies) != proxyCount {
		t.Fatalf("incorrect proxy count")
	}

	wg := sync.WaitGroup{}
	wg.Add(proxyCount * 2)

	for i := 0; i < proxyCount*2; i++ {
		go func() {
			list.SelectNextProxy()
			wg.Done()
		}()
	}

	wg.Wait()
}
