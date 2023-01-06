package proxies

import (
	"fmt"
	"testing"
)

func TestParseProxyNoCredentials(t *testing.T) {
	const host string = "127.0.0.1"
	const port string = "8080"

	proxy := Parse(fmt.Sprintf("%v:%v", host, port))

	if proxy.Host != host {
		t.Fatalf("incorrect proxy host - %v", proxy.Host)
	}

	if proxy.Port != port {
		t.Fatalf("incorrect proxy port - %v", proxy.Port)
	}
}

func TestParseProxyWithCredentials(t *testing.T) {
	const host string = "192.168.1.1"
	const port string = "1234"
	const username string = "admin"
	const password string = "password123"

	proxy := Parse(fmt.Sprintf("%v:%v:%v:%v", host, port, username, password))

	if proxy.Host != host {
		t.Fatalf("incorrect proxy host - %v", proxy.Host)
	}

	if proxy.Port != port {
		t.Fatalf("incorrect proxy port - %v", proxy.Port)
	}

	if proxy.Username != username {
		t.Fatalf("incorrect username - %v", proxy.Username)
	}

	if proxy.Password != password {
		t.Fatalf("incorrect password - %v", proxy.Password)
	}
}
