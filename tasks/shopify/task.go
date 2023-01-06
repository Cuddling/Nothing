package shopify

import (
	"Mystery/discord"
	"Mystery/profiles"
	"Mystery/proxies"
	"Mystery/tasks"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Task struct {
	tasks.Task
	Product               Product
	Variant               ProductVariant
	VariantInCart         int64
	CheckoutUrl           string
	CurrentPage           Page
	PreviousPage          Page
	ShippingRate          string
	PreviousQueueResponse PollQueueResponse
	SubmittedContactInfo  bool
}

// NewTaskShopify Returns a new Shopify task
func NewTaskShopify(site tasks.Website, profile *profiles.Profile, proxyList *proxies.ProxyList, mode tasks.TaskMode, inputs []string, sizes []string) Task {
	t := tasks.NewTask(site, profile, proxyList, mode, inputs, sizes)

	return Task{
		t,
		Product{},
		ProductVariant{},
		0,
		"",
		Page{},
		Page{},
		"",
		PollQueueResponse{},
		false,
	}
}

func (t *Task) Start() {
	if t.TaskManager != nil {
		t.TaskManager.TaskWaitGroup.Add(1)
	}

	go t.Run()
}

func (t *Task) Run() {
	t.SelectProxy()
	t.IsRunning = true

	for t.IsRunning {
		t.FlowMonitorProducts()
		t.FlowAddToCart(t.Variant.Id)
		t.FlowCreateCheckout()
		t.FlowPollQueue()

		// Checkpoint not support yet
		if t.CurrentPage.GetCheckoutStep() == CheckoutStepCheckpoint {
			t.UpdateStatus(&tasks.TaskStatus{
				Value: "Waiting For Checkpoint Captcha (Not Supported)",
				Level: tasks.StatusLevelError,
			}, true)

			t.Stop()
			return
		}

		// Login not supported yet
		if t.CurrentPage.GetCheckoutStep() == CheckoutStepLogin {
			t.UpdateStatus(&tasks.TaskStatus{
				Value: "Account Required (Not Supported)",
				Level: tasks.StatusLevelError,
			}, true)

			t.Stop()
			return
		}

		// Fully load the checkout page if we're on safe mode because anti-bot and other parameters need to be submitted if detected.
		if t.Mode == tasks.ModeShopifySafe {
			t.FlowLoadCheckoutPage()
		}

		// Below this point requires a valid checkout session.
		if t.CheckoutUrl == "" {
			continue
		}

		// On fast mode, immediately force submit contact info if it hasn't been submitted already. This allows us to skip the page loading
		// step which us a request - ultimately making fast mode faster.
		if t.Mode == tasks.ModeShopifyFast && !t.SubmittedContactInfo {
			t.FlowSubmitContactInfo(true)
		} else {
			t.FlowSubmitContactInfo(false)
		}

		// Shipping rates are required to be fetched now before they can be submitted. No way to avoid this for now that I can see.
		t.FlowFetchShippingRate()

		gateway := GetGateway(t.CurrentPage.GetShopId())

		// Go through the normal flow on safe mode or if we don't already have the gateway cached.
		if t.Mode == tasks.ModeShopifySafe || gateway == 0 {
			t.FlowSubmitShippingRate()
		}

		t.FlowCalculateTaxes()

		// Submits the payment and shipping rate in a single step if using fast mode, otherwise it
		// goes through the normal flow.
		if t.Mode == tasks.ModeShopifyFast && gateway != 0 {
			t.FlowSubmitPayment(true)
		} else {
			t.FlowSubmitPayment(false)
		}

		t.FlowProcessOrder()
		t.HandleCheckoutFailure()
		t.HandleCheckoutSuccess()
	}
}

// Stop the task from running
func (t *Task) Stop() {
	if t.TaskManager != nil {
		t.TaskManager.TaskWaitGroup.Done()
	}

	t.IsRunning = false
	t.Product = Product{}
	t.Variant = ProductVariant{}
	t.VariantInCart = 0
	t.CheckoutUrl = ""
	t.CurrentPage = Page{}
	t.PreviousPage = Page{}
	t.ShippingRate = ""
	t.PreviousQueueResponse = PollQueueResponse{}
	t.SubmittedContactInfo = false
	t.Log("Task Stopped")
}

// FlowMonitorProducts Monitors for products with error handling flow
func (t *Task) FlowMonitorProducts() {
	if !t.IsRunning {
		return
	}

	for !t.IsProductFound() {
		t.UpdateStatus(&tasks.TaskStatus{
			Value: "Monitoring",
			Level: tasks.StatusLevelInfo,
		}, true)

		resp, err := t.MonitorProducts()

		switch err {
		case nil:
			t.Log(fmt.Sprintf("Product Found: %v | %v (%v)", t.Product.Title, t.Variant.Option1, t.Variant.Id))
			t.ProductName = t.Product.Title
			t.ProductSize = t.Variant.Option1
		case ErrProductNotFound:
			t.UpdateStatus(&tasks.TaskStatus{
				Value: "Product Not Found",
				Level: tasks.StatusLevelInfo,
			}, true)

			time.Sleep(3500 * time.Millisecond)
		default:
			t.Log(err)
			t.UpdateStatus(&tasks.TaskStatus{
				Value: fmt.Sprintf("Error Monitoring Product (%v)", resp.StatusCode()),
				Level: tasks.StatusLevelError,
			}, false)

			time.Sleep(3500 * time.Millisecond)
		}
	}
}

// FlowAddToCart Adds item to the cart depending on the mode
func (t *Task) FlowAddToCart(variant int64) {
	if !t.IsRunning {
		return
	}

	// Not a valid variant or already in the cart
	if variant == 0 || t.VariantInCart == variant {
		return
	}

	t.UpdateStatus(&tasks.TaskStatus{
		Value: "Adding To Cart",
		Level: tasks.StatusLevelInfo,
	}, true)

	data, resp, err := t.AddToCart(variant)

	if err != nil {
		t.Log(err)
		t.UpdateStatus(&tasks.TaskStatus{
			Value: fmt.Sprintf("Error Adding To Cart (%v)", resp.StatusCode()),
			Level: tasks.StatusLevelError,
		}, false)

		time.Sleep(3500 * time.Millisecond)
		return
	}

	// When adding a variant to the cart, the response we get contains product information
	if t.Product.Id == 0 {
		t.Product = Product{
			Id:     data.ProductId,
			Title:  data.ProductTitle,
			Handle: data.Handle,
		}
	}

	t.Log(fmt.Sprintf("Item added to cart: %v (%v)", t.Product.Title, variant))
}

// FlowCreateCheckout Creates a checkout session for the flow
func (t *Task) FlowCreateCheckout() {
	if !t.IsRunning {
		return
	}

	if t.CheckoutUrl != "" || t.CurrentPage != (Page{}) {
		return
	}

	t.UpdateStatus(&tasks.TaskStatus{
		Value: "Creating Checkout",
		Level: tasks.StatusLevelInfo,
	}, true)

	url, resp, err := t.CreateCheckout()

	if err != nil {
		t.Log(err)
		t.UpdateStatus(&tasks.TaskStatus{
			Value: fmt.Sprintf("Error Creating Checkout (%v)", resp.StatusCode()),
			Level: tasks.StatusLevelError,
		}, false)

		time.Sleep(3500 * time.Millisecond)
		return
	}

	if url != "" {
		t.Log(fmt.Sprintf("Checkout session created: %v", url))
	} else {
		t.Log(fmt.Sprintf("Redirected to: %v", t.CurrentPage.Url))
	}
}

// FlowCreateCheckoutFast Creates a checkout session + adds to cart in one request for the flow
func (t *Task) FlowCreateCheckoutFast(variant int64) {
	if !t.IsRunning {
		return
	}

	if variant == 0 || t.VariantInCart == variant {
		return
	}

	if t.CheckoutUrl != "" || t.CurrentPage != (Page{}) {
		return
	}

	t.UpdateStatus(&tasks.TaskStatus{
		Value: fmt.Sprintf("Creating Checkout (Fast)"),
		Level: tasks.StatusLevelInfo,
	}, true)

	_, resp, err := t.CreateCheckoutFast(t.Variant.Id)

	if err != nil {
		t.Log(err)
		t.UpdateStatus(&tasks.TaskStatus{
			Value: fmt.Sprintf("Error Creating Checkout - Fast (%v)", resp.StatusCode()),
			Level: tasks.StatusLevelError,
		}, false)

		time.Sleep(3500 * time.Millisecond)
	}

	t.Log(fmt.Sprintf("Redirected to: %v", resp.RawResponse.Request.URL.String()))
}

// FlowLoadCheckoutPage Loads the initial checkout page
func (t *Task) FlowLoadCheckoutPage() {
	if !t.IsRunning {
		return
	}

	if t.CurrentPage != (Page{}) || t.CheckoutUrl == "" {
		return
	}

	t.UpdateStatus(&tasks.TaskStatus{
		Value: fmt.Sprintf("Fetching Checkout Page"),
		Level: tasks.StatusLevelInfo,
	}, true)

	page, resp, err := t.GetCheckoutPage(t.CheckoutUrl)

	if err != nil {
		t.Log(err)
		t.UpdateStatus(&tasks.TaskStatus{
			Value: fmt.Sprintf("Error Fetching Checkout Page (%v)", resp.StatusCode()),
			Level: tasks.StatusLevelError,
		}, false)

		time.Sleep(3500 * time.Millisecond)
		return
	}

	t.Log(fmt.Sprintf("Redirected to: %v", page.Url))
}

// FlowSubmitContactInfo Submits the contact info for the flow
func (t *Task) FlowSubmitContactInfo(force bool) {
	if !t.IsRunning {
		return
	}

	if !force && t.CurrentPage.GetCheckoutStep() != CheckoutStepContact {
		return
	}

	t.UpdateStatus(&tasks.TaskStatus{
		Value: fmt.Sprintf("Submitting Contact Info"),
		Level: tasks.StatusLevelInfo,
	}, true)

	page, resp, err := t.SubmitContactInfo()

	if err != nil {
		t.Log(err)
		t.UpdateStatus(&tasks.TaskStatus{
			Value: fmt.Sprintf("Error Submitting Contact Info (%v)", resp.StatusCode()),
			Level: tasks.StatusLevelError,
		}, false)

		time.Sleep(3500 * time.Millisecond)
		return
	}

	t.Log(fmt.Sprintf("Redirected to: %v", page.Url))
	t.SubmittedContactInfo = true
}

// FlowFetchShippingRate Fetches shipping rates for the flow
func (t *Task) FlowFetchShippingRate() {
	if !t.IsRunning {
		return
	}

	if t.CurrentPage.GetCheckoutStep() != CheckoutStepShippingMethod {
		return
	}

	for t.ShippingRate == "" {
		t.UpdateStatus(&tasks.TaskStatus{
			Value: fmt.Sprintf("Fetching Shipping Rate"),
			Level: tasks.StatusLevelInfo,
		}, true)

		rate, resp, err := t.FetchShippingRate()

		switch err {
		case nil:
			t.ShippingRate = rate
			t.Log(fmt.Sprintf("Fetched shipping rate: %v", t.ShippingRate))
		case ErrNoShippingRateAvailable:
			t.UpdateStatus(&tasks.TaskStatus{
				Value: fmt.Sprintf("No Shipping Rate Available"),
				Level: tasks.StatusLevelInfo,
			}, true)

			time.Sleep(1 * time.Second)
		default:
			t.Log(err)
			t.UpdateStatus(&tasks.TaskStatus{
				Value: fmt.Sprintf("Error Fetching Shipping Rate (%v)", resp.StatusCode()),
				Level: tasks.StatusLevelError,
			}, false)

			time.Sleep(3500 * time.Millisecond)
		}
	}
}

// FlowSubmitShippingRate Submits the shipping rate for the flow
func (t *Task) FlowSubmitShippingRate() {
	if !t.IsRunning {
		return
	}

	if t.CurrentPage.GetCheckoutStep() != CheckoutStepShippingMethod {
		return
	}

	if t.ShippingRate == "" {
		return
	}

	t.UpdateStatus(&tasks.TaskStatus{
		Value: fmt.Sprintf("Submitting Shipping Rate"),
		Level: tasks.StatusLevelInfo,
	}, true)

	page, resp, err := t.SubmitShippingRate(t.ShippingRate)

	if err != nil {
		t.Log(err)
		t.UpdateStatus(&tasks.TaskStatus{
			Value: fmt.Sprintf("Error Submitting Shipping Rate (%v)", resp.StatusCode()),
			Level: tasks.StatusLevelError,
		}, false)

		time.Sleep(3500 * time.Millisecond)
		return
	}

	t.Log(fmt.Sprintf("Redirected to: %v", page.Url))
}

// FlowCalculateTaxes Handles calculating taxes for the flow
func (t *Task) FlowCalculateTaxes() {
	if !t.IsRunning {
		return
	}

	for t.CurrentPage.GetCheckoutStep() == CheckoutStepCalculatingTaxes {
		t.UpdateStatus(&tasks.TaskStatus{
			Value: fmt.Sprintf("Calculating Taxes"),
			Level: tasks.StatusLevelInfo,
		}, true)

		done, resp, err := t.CalculateTaxes()

		if err != nil {
			t.Log(err)
			t.UpdateStatus(&tasks.TaskStatus{
				Value: fmt.Sprintf("Error Calculating Taxes (%v)", resp.StatusCode()),
				Level: tasks.StatusLevelError,
			}, false)

			time.Sleep(3500 * time.Millisecond)
		}

		if !done {
			t.Log("Still calculating taxes")
			time.Sleep(1 * time.Millisecond)
			continue
		}

		t.Log("Finished calculating taxes!")
		t.Log(fmt.Sprintf("Redirected to: %v", t.CurrentPage.Url))
	}
}

// FlowSubmitPayment Submits payment for the flow
func (t *Task) FlowSubmitPayment(fast bool) {
	if !t.IsRunning {
		return
	}

	if fast {
		gateway := GetGateway(t.CurrentPage.GetShopId())

		if t.ShippingRate == "" || gateway == 0 {
			return
		}
	} else {
		if t.CurrentPage.GetCheckoutStep() != CheckoutStepPaymentMethod {
			return
		}
	}

	t.UpdateStatus(&tasks.TaskStatus{
		Value: fmt.Sprintf("Fetching Payment Token"),
		Level: tasks.StatusLevelInfo,
	}, true)

	token, resp, err := t.GetPaymentToken()

	if err != nil {
		t.Log(err)
		t.UpdateStatus(&tasks.TaskStatus{
			Value: fmt.Sprintf("Error Fetching Payment Token (%v)", resp.StatusCode()),
			Level: tasks.StatusLevelError,
		}, false)

		time.Sleep(3500 * time.Millisecond)
		return
	}

	t.UpdateStatus(&tasks.TaskStatus{
		Value: fmt.Sprintf("Submitting Payment"),
		Level: tasks.StatusLevelInfo,
	}, true)

	page, resp, err := t.SubmitPayment(fast, token, t.ShippingRate)

	if err != nil {
		t.Log(err)
		t.UpdateStatus(&tasks.TaskStatus{
			Value: fmt.Sprintf("Error Submitting Payment (%v)", resp.StatusCode()),
			Level: tasks.StatusLevelError,
		}, false)

		time.Sleep(3500 * time.Millisecond)
		return
	}

	t.Log(fmt.Sprintf("Redirected to : %v", page.Url))
}

// FlowProcessOrder Polls the order processing page for the flow
func (t *Task) FlowProcessOrder() {
	if !t.IsRunning {
		return
	}

	for t.CurrentPage.GetCheckoutStep() == CheckoutStepProcessing {
		t.UpdateStatus(&tasks.TaskStatus{
			Value: fmt.Sprintf("Processing Order"),
			Level: tasks.StatusLevelImportant,
		}, true)

		resp, err := t.ProcessOrder()

		if err != nil {
			t.Log(err)
			t.UpdateStatus(&tasks.TaskStatus{
				Value: fmt.Sprintf("Error Processing Order (%v)", resp.StatusCode()),
				Level: tasks.StatusLevelError,
			}, false)

			time.Sleep(3500 * time.Millisecond)
			continue
		}

		t.Log(fmt.Sprintf("Redirected to: %v", t.CurrentPage.Url))

		// Add a delay before polling the processing step again
		if t.CurrentPage.GetCheckoutStep() == CheckoutStepProcessing {
			time.Sleep(1 * time.Second)
			continue
		}
	}
}

// FlowPollQueue Handles polling the checkout queue
func (t *Task) FlowPollQueue() {
	if !t.IsRunning {
		return
	}

	for t.CurrentPage.GetCheckoutStep() == CheckoutStepQueue {
		// Wait until we're allowed to poll the queue again or if it has never been polled.
		if t.PreviousQueueResponse != (PollQueueResponse{}) && time.Now().UTC().Before(t.PreviousQueueResponse.Data.Poll.PollAfter) {
			continue
		}

		t.UpdateStatus(&tasks.TaskStatus{
			Value: fmt.Sprintf("Polling Queue"),
			Level: tasks.StatusLevelInfo,
		}, true)

		data, _, err := t.PollQueue()

		if err != nil {
			t.Log(err)
			time.Sleep(3500 * time.Millisecond)
			continue
		}

		t.PreviousQueueResponse = data

		switch t.PreviousQueueResponse.Data.Poll.GetTypename() {
		case QueuePollContinue:
			t.Log(fmt.Sprintf("Still In Queue | ETA: %v | Poll After: %v", data.Data.Poll.QueueEtaSeconds, data.Data.Poll.PollAfter))
			time.Sleep(3000 * time.Millisecond)
		case QueuePollComplete:
			t.Log("Queue completed!")

			url := t.PreviousPage.Url

			// Fallback to the checkout URL
			if url == "" {
				url = t.CheckoutUrl
			}

			// Final fallback. Clear the previous and current pages, so the checkout session can finally be created.
			if url == "" {
				t.CurrentPage = Page{}
				t.PreviousPage = Page{}
				return
			}

			if !t.IsRunning {
				return
			}

			t.UpdateStatus(&tasks.TaskStatus{
				Value: fmt.Sprintf("Fetching Checkout Page"),
				Level: tasks.StatusLevelInfo,
			}, true)

			t.Log(url)
			page, resp, err := t.GetCheckoutPage(url)

			if err != nil {
				t.Log(err)
				t.UpdateStatus(&tasks.TaskStatus{
					Value: fmt.Sprintf("Error Fetching Checkout Page (%v)", resp.StatusCode()),
					Level: tasks.StatusLevelError,
				}, false)

				time.Sleep(3500 * time.Millisecond)
				continue
			}

			// Reset the queue response since we're no longer in queue.
			t.PreviousQueueResponse = PollQueueResponse{}
			t.UpdateCurrentPage(page)
		default:
			t.Log(fmt.Sprintf("Unknown queue poll response: %v", data))
		}
	}
}

// HandleCheckoutFailure Handles when the checkout has failed in some way
func (t *Task) HandleCheckoutFailure() {
	if t.CurrentPage.GetCheckoutStep() != CheckoutStepPaymentMethod {
		return
	}

	notice, err := t.CurrentPage.GetNoticeText()

	// Don't think this will ever be hit, but in the event that there is no notice text,
	// it'll resort to using the raw html.
	if err != nil {
		notice = t.CurrentPage.Html
	}

	// Out of stock. Task doesn't stop here, as we'll still be running for restocks
	if strings.Contains(t.CurrentPage.Url, "stock_problems") {
		t.UpdateStatus(&tasks.TaskStatus{
			Value: fmt.Sprintf("Payment Failed - Out of Stock"),
			Level: tasks.StatusLevelInfo,
		}, true)

		time.Sleep(3500 * time.Millisecond)
		return
	}

	// With fast mode, this can happen when the total price is incorrect
	if strings.Contains(strings.ToLower(notice), "order total has changed") {
		t.UpdateStatus(&tasks.TaskStatus{
			Value: fmt.Sprintf("Payment Failed - Incorrect Order Total"),
			Level: tasks.StatusLevelError,
		}, true)
		return
	}

	// Handle common payment decline reasons
	for _, reason := range []string{"issue processing", "problem processing", "declined", "insufficient"} {
		if strings.Contains(strings.ToLower(notice), reason) {
			t.Log(fmt.Sprintf("Payment Declined - %v", notice))
			t.UpdateStatus(&tasks.TaskStatus{
				Value: "Payment Declined",
				Level: tasks.StatusLevelError,
			}, false)

			t.SendWebhook(false)
			t.Stop()
			return
		}
	}

	t.Log(fmt.Sprintf("Payment Declined - %v", notice))
	t.UpdateStatus(&tasks.TaskStatus{
		Value: "Payment Error",
		Level: tasks.StatusLevelError,
	}, false)

	t.SendWebhook(false)
	t.Stop()
}

// HandleCheckoutSuccess Handles when the checkout has succeeded
func (t *Task) HandleCheckoutSuccess() {
	if t.CurrentPage.GetCheckoutStep() != CheckoutStepOrderConfirmation {
		return
	}

	t.UpdateStatus(&tasks.TaskStatus{
		Value: "Checked Out!",
		Level: tasks.StatusLevelSuccess,
	}, true)

	t.SendWebhook(true)
	t.Stop()
}

// IsProductFound Returns if the product is found
func (t *Task) IsProductFound() bool {
	return t.Variant != (ProductVariant{})
}

// UpdateCurrentPage Updates the current and previous pages of the task
func (t *Task) UpdateCurrentPage(p Page) {
	t.PreviousPage = t.CurrentPage
	t.CurrentPage = p
}

// SendWebhook Sends a webhook to Discord
func (t *Task) SendWebhook(success bool) {
	failureReason := "None"
	proxyList := "None"
	orderNumber := "None"
	orderLink := ""

	if t.ProxyList != nil {
		proxyList = t.ProxyList.Name
	}

	if success {
		orderNumber = t.CurrentPage.GetOrderNumber()
		orderLink = t.CurrentPage.Url
	} else {
		failureReason, _ = t.CurrentPage.GetNoticeText()
	}

	discord.SendCheckoutWebhook(discord.CheckoutWebhookData{
		Success:       success,
		FailureReason: failureReason,
		Site:          t.Site.Name,
		Mode:          tasks.ModeToString(t.Mode),
		ProductTitle:  t.CurrentPage.GetProductTitle(),
		ProductSize:   t.CurrentPage.GetProductSize(),
		ProductImage:  t.CurrentPage.GetProductImage(),
		Profile:       t.Profile.Name,
		ProxyList:     proxyList,
		Email:         t.Profile.ShippingAddress.Email,
		OrderNumber:   orderNumber,
		OrderLink:     orderLink,
		Quantity:      strconv.Itoa(t.Quantity),
	})
}
