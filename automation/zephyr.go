package automation

import (
	"encoding/json"
	"fmt"
	"github.com/sacOO7/gowebsocket"
	"log"
)

// ZephyrMonitorMessage Standard monitor message which only parses the type, so the correct struct can be parsed.
type ZephyrMonitorMessage struct {
	Type string `json:"type"`
	Data json.RawMessage
}

// ZephyrMonitorSiteList The list of sites that is being monitored.
var ZephyrMonitorSiteList *MonitorMessagePinConfig

// ZephyrMonitorCurrentAntibot The current anti-bot status for all websites.
var ZephyrMonitorCurrentAntibot *ZephyrMonitorMessageAntibot

// ZephyrMonitorPreviousAntibot The previous anti-bot status for all websites.
var ZephyrMonitorPreviousAntibot *ZephyrMonitorMessageAntibot

// ZephyrMonitorChannel Channel for zephyr monitor products
var ZephyrMonitorChannel = make(chan *ZephyrMonitorLive, 1)

// ConnectToZephyrMonitor Connects to the Z-AIO monitor.
func ConnectToZephyrMonitor() {
	const MonitorKey string = ""
	const MonitorUri string = ""

	socket := gowebsocket.New(MonitorUri)

	socket.OnConnected = func(socket gowebsocket.Socket) {
		log.Println("Connected to monitor. Authenticating...")
		socket.SendText(fmt.Sprintf("{\"server\":null,\"type\":\"auth\",\"version\":\"1.9.76\",\"key\":\"%v\"}", MonitorKey))
		socket.SendText(fmt.Sprintf("{\"wsServerDisconnections\":0,\"shippingRates\":[],\"tasks247Footsites\":0,\"interval\":3600000,\"type\":\"tasksStats\",\"version\":\"1.9.76\",\"tasks247\":0,\"key\":\"%v\",\"tasks\":0,\"wsClientDisconnections\":0}", MonitorKey))
		log.Println("Sent authentication details.")
	}

	socket.OnConnectError = func(err error, socket gowebsocket.Socket) {
		log.Println("Monitor Connection Error: ", err)
		ConnectToZephyrMonitor()
		return
	}

	socket.OnTextMessage = func(message string, socket gowebsocket.Socket) {
		handleZephyrMonitorMessage(message)
	}

	socket.OnBinaryMessage = func(data []byte, socket gowebsocket.Socket) {
		log.Println("Received Binary Message: ", data)
	}

	socket.OnPingReceived = func(data string, socket gowebsocket.Socket) {
		socket.SendText(fmt.Sprintf("{\"type\":\"ping\",\"key\":\"%v\"}", MonitorKey))
	}

	socket.OnPongReceived = func(data string, socket gowebsocket.Socket) {
		log.Println("Received Pong: " + data)
	}

	socket.OnDisconnected = func(err error, socket gowebsocket.Socket) {
		log.Println("Disconnected from monitor. Reconnecting.")
		ConnectToZephyrMonitor()
		return
	}

	log.Println("Connecting to monitor...")
	socket.Connect()
}

// HandleMonitorMessage Handles incoming monitor messages.
func handleZephyrMonitorMessage(message string) {
	var rawMsg ZephyrMonitorMessage

	if err := json.Unmarshal([]byte(message), &rawMsg); err != nil {
		log.Println("Failed to unmarshal monitor message.")
		log.Println(err)
		return
	}

	switch rawMsg.Type {
	case "pinConfig":
		if err := json.Unmarshal([]byte(message), &ZephyrMonitorSiteList); err != nil {
			log.Printf("Failed to unmarshal monitor message: %v\n", err)
			return
		}

		log.Printf("Received %v sites from the monitor.\n", len(ZephyrMonitorSiteList.Body))
	case "livemonitorInit":
		return
	case "livemonitorInitFootsites":
		return
	case "livemonitorInitSupreme":
		return
	case "logKeywords":
		return
	case "botFeatures":
		return
	case "shopifyAntibot":
		ZephyrMonitorPreviousAntibot = ZephyrMonitorCurrentAntibot

		if err := json.Unmarshal([]byte(message), &ZephyrMonitorCurrentAntibot); err != nil {
			log.Printf("Failed to unmarshal monitor message: %v\n", err)
			return
		}

		log.Printf("Received anti-bot data for %v sites.\n", len(ZephyrMonitorCurrentAntibot.Sites))
	case "livemonitor":
		log.Printf(message)

		var liveProduct ZephyrMonitorLive

		if err := json.Unmarshal([]byte(message), &liveProduct); err != nil {
			log.Printf("Failed to unmarshal monitor message: %v\n", err)
			return
		}

		// Only caring about shopify products for now.
		if liveProduct.Body.Type != "shopify" {
			return
		}

		ZephyrMonitorChannel <- &liveProduct
		p := liveProduct.Body.Payload.Product
		log.Printf("Received incoming liveProduct - Store: %v | Title: %v | Price: %v\n", liveProduct.Body.Payload.Store, p.Title, p.Variants[0].Price)
	case "pong":
		return
	}
}
