package tasks

import (
	"Mystery/profiles"
	"Mystery/proxies"
	"Mystery/utils"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
	"net/url"
	"time"
)

type Task struct {
	TaskManager   *Manager
	Group         *TaskGroup
	Id            string
	Status        *TaskStatus
	IsRunning     bool
	Client        *resty.Client
	Site          Website
	Profile       *profiles.Profile
	ProxyList     *proxies.ProxyList
	Mode          TaskMode
	MonitorInputs []string
	Sizes         []string
	Quantity      int
	ProductName   string
	ProductSize   string
}

type TaskStatus struct {
	Value string
	Level TaskStatusLevel
}

type TaskStatusLevel int

const (
	StatusLevelInfo TaskStatusLevel = iota
	StatusLevelImportant
	StatusLevelError
	StatusLevelSuccess
)

type TaskRunner interface {
	Start()
	Stop()
}

// NewTask Creates and returns a new task
func NewTask(site Website, profile *profiles.Profile, proxyList *proxies.ProxyList, mode TaskMode, inputs []string, sizes []string) Task {
	t := Task{
		Id:     utils.GenerateRandomString(8),
		Client: resty.New(),
		Status: &TaskStatus{
			Value: "Idle",
			Level: StatusLevelInfo,
		},
		Site:          site,
		Profile:       profile,
		ProxyList:     proxyList,
		Mode:          mode,
		MonitorInputs: inputs,
		Sizes:         sizes,
		Quantity:      1,
	}

	t.Client.SetTimeout(1 * time.Minute)
	return t
}

// SelectProxy Selects the proxy that the task will use to check out
func (t *Task) SelectProxy() {
	if t.ProxyList == nil {
		t.Log("Using Proxy: Localhost")
		return
	}

	t.Client.RemoveProxy()

	proxy := t.ProxyList.SelectNextProxy()

	urlStr := ""

	if proxy.Username != "" && proxy.Password != "" {
		urlStr = fmt.Sprintf("http://%v:%v@%v:%v", proxy.Username, proxy.Password, proxy.Host, proxy.Port)
	} else {
		urlStr = fmt.Sprintf("http://%v:%v", proxy.Host, proxy.Port)
	}

	uri, _ := url.Parse(urlStr)

	t.Client.SetTransport(&http.Transport{
		Proxy: http.ProxyURL(uri),
	})

	t.Client.SetProxy(urlStr)
	t.Log(fmt.Sprintf("Using Proxy: %v:%v", proxy.Host, proxy.Port))
}

// UpdateStatus Updates the status of the task (ex. Idle / Adding To Cart / Processing)
func (t *Task) UpdateStatus(s *TaskStatus, log bool) {
	t.Status = s

	if log {
		t.Log(s.Value)
	}
}

// Log a string to stdout / file
func (t *Task) Log(l interface{}) {
	now := time.Now()
	fmt.Printf("[%d:%d:%d] [TASK: %v] %v\n", now.Hour(), now.Minute(), now.Second(), t.Id, l)
}
