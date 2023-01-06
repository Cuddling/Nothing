package tasks

import (
	"Mystery/automation"
	"sync"
	"testing"
)

func TestTaskGroups(t *testing.T) {
	m := NewManager()

	group := NewTaskGroup("Test Group")

	for i := 0; i < 3; i++ {
		task := NewTask(Website{}, nil, nil, ModeShopifyFast, []string{}, []string{})
		group.AddTask(&task)
	}

	m.AddTaskGroup(&group)

	if len(m.Tasks) != 3 {
		t.Fatalf("Incorrect task count")
	}

	if len(m.TaskGroups) != 1 {
		t.Fatalf("Incorrect group count")
	}

	m.RemoveTaskGroup(&group)

	if len(m.TaskGroups) != 0 {
		t.Fatalf("Incorrect group count (removal)")
	}

	if len(m.Tasks) != 0 {
		t.Fatalf("Incorrect task count (removal)")
	}
}

func TestAutomationChannel(t *testing.T) {
	m := NewManager()

	m.Automations = append(m.Automations, &automation.Automation{
		Name:             "Automation Test (+a)",
		MonitorInputs:    []string{"+a"},
		Sizes:            []string{},
		Profiles:         nil,
		ProxyList:        "",
		CheckUrl:         true,
		PriceMinimum:     0,
		PriceMaximum:     0,
		Quantity:         1,
		TotalTaskCount:   50,
		SiteWhitelist:    nil,
		SiteBlacklist:    nil,
		PaymentRetries:   0,
		StopAfterMinutes: 5,
	})

	m.Automations = append(m.Automations, &automation.Automation{
		Name:             "Automation Test (+dunk)",
		MonitorInputs:    []string{"+dunk"},
		Sizes:            []string{},
		Profiles:         nil,
		ProxyList:        "",
		CheckUrl:         true,
		PriceMinimum:     0,
		PriceMaximum:     0,
		Quantity:         1,
		TotalTaskCount:   50,
		SiteWhitelist:    nil,
		SiteBlacklist:    nil,
		PaymentRetries:   0,
		StopAfterMinutes: 5,
	})

	m.HandleNewAutomationProducts()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go automation.ConnectToZephyrMonitor()
	wg.Wait()
}
