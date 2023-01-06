package discord

import (
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/webhook"
	"github.com/disgoorg/snowflake/v2"
	"time"
)

type CheckoutWebhookData struct {
	Success       bool
	FailureReason string
	Site          string
	Mode          string
	ProductTitle  string
	ProductSize   string
	ProductImage  string
	Profile       string
	ProxyList     string
	Email         string
	OrderNumber   string
	OrderLink     string
	Quantity      string
}

// SendCheckoutWebhook Sends a success/failure checkout webhook to Discord
func SendCheckoutWebhook(data CheckoutWebhookData) {
	client := webhook.New(snowflake.ID(0), "")
	embed := discord.NewEmbedBuilder()

	if data.Success {
		embed.SetAuthor("Successful Checkout! ✅", "", "https://i.imgur.com/RLauXcd.png")
		embed.SetColor(0x000001)
	} else {
		embed.SetAuthor("Checkout Failed! ❌", "", "https://i.imgur.com/RLauXcd.png")
		embed.SetDescriptionf("**Reason:** %v", data.FailureReason)
		embed.SetColor(0xFF0000)
	}

	if data.ProductImage != "" {
		embed.SetThumbnail(data.ProductImage)
	}

	embed.AddField("Site", data.Site, true)
	embed.AddField("Mode", data.Mode, true)
	embed.AddField("Product", data.ProductTitle, true)
	embed.AddField("Size", data.ProductSize, true)
	embed.AddField("Quantity", data.Quantity, true)
	embed.AddField("Profile", fmt.Sprintf("||%v||", data.Profile), true)
	embed.AddField("Email", fmt.Sprintf("||%v||", data.Email), true)
	embed.AddField("Proxy List", fmt.Sprintf("||%v||", data.ProxyList), true)

	if data.OrderLink != "" {
		embed.AddField("Order Number", fmt.Sprintf("||[%v](%v)||", data.OrderNumber, data.OrderLink), true)
	} else if data.OrderNumber != "" {
		embed.AddField("Order Number", fmt.Sprintf("||%v||", data.OrderNumber), true)
	} else {
		embed.AddField("Order Number", "None", true)
	}

	embed.SetFooter("Nothing", "https://i.imgur.com/RLauXcd.png")
	embed.SetTimestamp(time.Now())

	_, err := client.CreateEmbeds([]discord.Embed{embed.Build()})

	if err != nil {
		fmt.Println(err)
		return
	}
}
