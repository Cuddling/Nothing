package automation

import (
	"log"
)

type MonitorMessagePinConfig struct {
	Type string
	Body map[string]interface{}
}

func (m *MonitorMessagePinConfig) PrintSites() {
	for k, _ := range m.Body {
		log.Println(k)
	}
}
