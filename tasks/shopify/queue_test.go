package shopify

import (
	"Mystery/tasks"
	"testing"
)

// Tests polling the queue on a website
func TestPollQueue(t *testing.T) {
	site := tasks.Website{Name: "Oneness Boutique", Url: "https://www.onenessboutique.com"}
	variant := 39885244498004
	task := NewTaskShopify(site, &testProfile, nil, tasks.ModeShopifySafe, []string{string(rune(variant))}, []string{"L"})
	_, _, _ = task.CreateCheckout()
	data, _, _ := task.PollQueue()

	if data.Data.Poll.GetTypename() == QueuePollUnknown {
		t.Fatalf("expected a known queue typename. Got %v", data.Data.Poll.GetTypename())
	}
}
