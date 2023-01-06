package automation

import "log"

type ZephyrMonitorMessageAntibot struct {
	Type  string                            `json:"type"`
	Sites []ZephyrMonitorMessageAntibotSite `json:"sites"`
}

type ZephyrMonitorMessageAntibotSite struct {
	Antibot bool   `json:"antibot"`
	Uuid    string `json:"uuid"`
	Url     string `json:"url"`
}

func (m *ZephyrMonitorMessageAntibot) PrintSites() {
	for _, s := range m.Sites {
		log.Printf("%v - %v\n", s.Url, s.Antibot)
	}
}
