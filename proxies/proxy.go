package proxies

import "strings"

type Proxy struct {
	Host     string
	Port     string
	Username string
	Password string
	Raw      string
}

// Parse parses raw proxy into a proxy struct object
func Parse(str string) Proxy {
	split := strings.Split(str, ":")

	// At minimum a proxy needs to have a host and port
	if len(split) < 2 {
		return Proxy{}
	}

	proxy := Proxy{Raw: str}

	for i, val := range split {
		switch i {
		case 0:
			proxy.Host = val
		case 1:
			proxy.Port = val
		case 2:
			proxy.Username = val
		case 3:
			proxy.Password = val
		}
	}

	return proxy
}
