package automation

import (
	"Mystery/utils"
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/webhook"
	"github.com/disgoorg/log"
	"github.com/disgoorg/snowflake/v2"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type Automation struct {
	Name             string   // A unique name for the automation group.
	MonitorInputs    []string // List of monitor inputs whether it be keywords, links, or variants.
	Sizes            []string // The size range it will check out. Leave blank for random.
	Profiles         []string // List of profiles it will use to check out. It distributes profiles equally with TotalTaskCount.
	ProxyList        string   // The proxy list it will use to check out.
	CheckUrl         bool     // If it checks the url/handle of the item
	PriceMinimum     int      // The minimum price in USD the product needs to fall within.
	PriceMaximum     int      // The maximum price in USD the product needs to fall within.
	Quantity         int      // The amount of the given product it will attempt to check out.
	TotalTaskCount   int      // The amount of total tasks that will be run.
	SiteWhitelist    []string // If not empty, only go for these specific sites and ignore SiteBlacklist. Otherwise, it uses all sites.
	SiteBlacklist    []string // Ignores products from specific sites.
	PaymentRetries   int      // The amount of times it'll attempt to retry payment submission.
	StopAfterMinutes int      // The time in minutes the tasks will stop after.
}

// IsProductMatch Returns if the product is a match for the automation
func (a *Automation) IsProductMatch(live *ZephyrMonitorLive) bool {
	if len(a.MonitorInputs) == 0 {
		return false
	}

	product := live.Body.Payload.Product

	for _, input := range a.MonitorInputs {
		switch utils.GetMonitorInputTypeFromString(input) {
		case utils.MonitorInputTypeKeywords:
			if utils.IsKeywordMatch(product.Title, input) {
				return true
			}

			if a.CheckUrl && utils.IsKeywordMatch(product.Handle, input) {
				return true
			}

			if len(product.Variants) > 0 {
				variant := product.Variants[0]

				// Check SKU
				if utils.IsKeywordMatch(variant.Sku, input) {
					return true
				}

				// Check variant name as a last resort
				if utils.IsKeywordMatch(variant.Name, input) {
					return true
				}
			}

		}
	}

	return false
}

// IsWebsiteMatch Checks if the website is one that is being checked out for automation
func (a *Automation) IsWebsiteMatch(live *ZephyrMonitorLive) bool {
	productStoreUrl, err := url.Parse(live.Body.Payload.Store)

	if err != nil {
		log.Error(err)
		return false
	}

	// Whitelist takes precedence over blacklist. Both cannot be used at the same time.
	// Check for whitelist first
	if len(a.SiteWhitelist) > 0 {
		for _, site := range a.SiteWhitelist {
			siteUrl, err := url.Parse(site)

			if err != nil {
				log.Error(err)
				return false
			}

			if siteUrl.Hostname() == productStoreUrl.Hostname() {
				return true
			}
		}
	} else if len(a.SiteBlacklist) > 0 {
		for _, site := range a.SiteBlacklist {
			siteUrl, err := url.Parse(site)

			if err != nil {
				log.Error(err)
				return false
			}

			// The site for this product is blacklisted
			if siteUrl.Hostname() == productStoreUrl.Hostname() {
				return false
			}
		}
	}

	return true
}

// IsPriceMatch Checks if the price of the product is within range
func (a *Automation) IsPriceMatch(live *ZephyrMonitorLive) bool {
	variants := live.Body.Payload.Product.Variants

	if len(variants) == 0 {
		log.Error("No variants found for product")
		return false
	}

	price, err := variants[0].Price.Float64()

	if err != nil {
		log.Error(err)
		return false
	}

	return price >= float64(a.PriceMinimum) && price <= float64(a.PriceMaximum)
}

// GetMatchingSizeVariants Gets variants that match the specified sizes
func (a *Automation) GetMatchingSizeVariants(live *ZephyrMonitorLive) []ZephyrMonitorLiveProductVariant {
	productVariants := live.Body.Payload.Product.Variants
	var matched []ZephyrMonitorLiveProductVariant

	for _, variant := range productVariants {
		if !variant.Available {
			continue
		}

		if len(a.Sizes) == 0 {
			matched = append(matched, variant)
		} else {
			for _, size := range a.Sizes {
				// Shoe sizes are determined if the variant option contains ANY number
				if regexp.MustCompile(`\d`).MatchString(variant.Title) {
					// Isolating the numbers from the size, so we can check against that
					regex := regexp.MustCompile("[^0-9.]+")
					rawSize := regex.ReplaceAllString(variant.Title, "")

					if rawSize == size {
						matched = append(matched, variant)
					}
				} else {
					// Otherwise, we can safely assume it's clothing and match the exact size.
					// Not using contains because "L" can pick up XL, XXL, etc.
					if variant.Title == size {
						matched = append(matched, variant)
					}
				}
			}
		}
	}

	return matched
}

// SendWebhook SendAutomationWebhook Sends a webhook to signal that automation for a product has started/stopped
func (a *Automation) SendWebhook(started bool, live *ZephyrMonitorLive) {
	client := webhook.New(snowflake.ID(0), "")
	embed := discord.NewEmbedBuilder()

	if started {
		embed.SetAuthor("\U0001F7E2 Automation Started", "", "https://i.imgur.com/RLauXcd.png")
		embed.SetColor(0x00FF00)
	} else {
		embed.SetAuthor("🔴 Automation Stopped", "", "https://i.imgur.com/RLauXcd.png")
		embed.SetColor(0xFF0000)
	}

	store := live.Body.Payload.Store
	product := live.Body.Payload.Product

	embed.SetThumbnail(product.GetImage())
	embed.AddField("Automation Name", a.Name, false)
	embed.AddField("Site", store, false)
	embed.AddField("Product", fmt.Sprintf("[%v](%vproducts/%v)", product.Title, store, product.Handle), false)

	sizes := "Random"

	if len(a.Sizes) > 0 {
		sizes = strings.Join(a.Sizes, ", ")
	}

	embed.AddField("Sizes", sizes, true)
	embed.AddField("Quanitty", fmt.Sprintf("%v", a.Quantity), true)
	embed.AddField("Task Count", fmt.Sprintf("%v", a.TotalTaskCount), true)
	embed.AddField("Payment Retries", fmt.Sprintf("%v", a.PaymentRetries), true)
	embed.AddField("Stop After Minutes", fmt.Sprintf("%v", a.StopAfterMinutes), true)

	embed.SetFooter("Nothing", "https://i.imgur.com/RLauXcd.png")
	embed.SetTimestamp(time.Now())

	_, err := client.CreateEmbeds([]discord.Embed{embed.Build()})

	if err != nil {
		fmt.Println(err)
		return
	}
}
