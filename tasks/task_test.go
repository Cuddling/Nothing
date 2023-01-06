package tasks

import (
	"Mystery/proxies"
	"fmt"
	"testing"
)

func TestTaskWithProxy(t *testing.T) {
	list := proxies.NewProxyList("Test", "127.0.0.1:8080:username:password")
	task := NewTask(Website{}, nil, &list, ModeShopifySafe, nil, nil)
	task.SelectProxy()

	resp, err := task.Client.R().Get("https://api.ipify.org?format=json")

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(resp.Body()))
}
