package shopify

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
	"net/url"
	"time"
)

type PollQueueResponse struct {
	Data PollQueueResponseData `json:"data,omitempty"`
}

type PollQueueResponseData struct {
	Poll PollQueueResponsePoll `json:"poll,omitempty"`
}

type PollQueueResponsePoll struct {
	Typename        string    `json:"__typename,omitempty"`
	QueueEtaSeconds int       `json:"queueEtaSeconds,omitempty"`
	Token           string    `json:"token,omitempty"`
	PollAfter       time.Time `json:"pollAfter,omitempty"`
}

type PollQueueTypename int

const (
	QueuePollContinue PollQueueTypename = iota
	QueuePollComplete
	QueuePollUnknown
)

// GetTypename Returns the queue polling typename in a better format
func (p *PollQueueResponsePoll) GetTypename() PollQueueTypename {
	switch p.Typename {
	case "PollContinue":
		return QueuePollContinue
	case "PollComplete":
		return QueuePollComplete
	default:
		return QueuePollUnknown
	}
}

// PollQueue Polls the checkout queue
func (t *Task) PollQueue() (PollQueueResponse, *resty.Response, error) {
	cookie := t.GetCheckoutQueueTokenCookie()

	if cookie == nil {
		return PollQueueResponse{}, nil, errors.New("no checkout queue token cookie found")
	}

	resp, err := t.Client.R().
		SetHeader("Accept", "*/*").
		SetHeader("Accept-Encoding", "gzip, deflate").
		SetHeader("Accept-Language", "en-US,en.q=0.9").
		SetHeader("Connection", "keep-alive").
		SetHeader("Content-Type", "application/json").
		SetHeader("Host", t.Site.GetHostName()).
		SetHeader("Referer", fmt.Sprintf("%v/throttle/queue", t.Site.Url)).
		SetHeader("sec-ch-ua", "\" Not A.Brand\".v=\"99\", \"Chromium\".v=\"102\", \"Google Chrome\".v=\"102\"").
		SetHeader("sec-ch-ua-mobile", "?0").
		SetHeader("sec-ch-ua-platform", "\"Windows\"").
		SetHeader("Sec-Fetch-Dest", "empty").
		SetHeader("Sec-Fetch-Mode", "cors").
		SetHeader("Sec-Fetch-Site", "same-origin").
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0. Win64. x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.62 Safari/537.36").
		SetBody(map[string]interface{}{
			"query": "\n      {\n        poll(token: $token) {\n          token\n          pollAfter\n          queueEtaSeconds\n          productVariantAvailability {\n            id\n            available\n          }\n        }\n      }\n    ",
			"variables": map[string]interface{}{
				"token": cookie.Value,
			},
		}).
		Post(fmt.Sprintf("%v/queue/poll", t.Site.Url))

	if err != nil {
		return PollQueueResponse{}, nil, err
	}

	if resp.IsError() {
		return PollQueueResponse{}, resp, errors.New(fmt.Sprintf("poll queue failed (%v)", resp.StatusCode()))
	}

	var data PollQueueResponse
	err = json.Unmarshal(resp.Body(), &data)

	if err != nil {
		return PollQueueResponse{}, resp, err
	}

	// Set the new cookie so that it can be used again on the next request.
	// This isn't being done automatically by the server for some reason.
	if data.Data.Poll.Token != "" {
		cookie.Value = data.Data.Poll.Token
		siteUrl, _ := url.Parse(t.Site.Url)
		t.Client.GetClient().Jar.SetCookies(siteUrl, []*http.Cookie{cookie})
	}

	return data, resp, nil
}

// GetCheckoutQueueTokenCookie Returns the checkout queue token cookie if one exists
func (t *Task) GetCheckoutQueueTokenCookie() *http.Cookie {
	siteUrl, _ := url.Parse(t.Site.Url)

	for _, cookie := range t.Client.GetClient().Jar.Cookies(siteUrl) {
		if cookie.Name != "_checkout_queue_token" {
			continue
		}

		return cookie
	}

	return nil
}
