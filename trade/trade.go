package trade

import (
	"EX_okexquant/config"
	"fmt"
)

var OKexClient *Client

func InitTrade() {
	OKexClient = newOKExClient()
	fmt.Println("[InitTrade] trade success.")
}

func newOKExClient() *Client {
	var con Config
	con.Endpoint = config.Config.Trade.Endpoint
	con.WSEndpoint = config.Config.Trade.WSEndpoint

	con.ApiKey = config.Config.Trade.ApiKey
	con.SecretKey = config.Config.Trade.SecretKey
	con.Passphrase = config.Config.Trade.Passphrase
	con.TimeoutSecond = config.Config.Trade.TimeoutSecond
	con.IsPrint = config.Config.Trade.IsPrint
	con.I18n = config.Config.Trade.I18n

	client := NewClient(con)
	return client
}

func NewClientByParam(apiKey, secretKey, passphrase string) *Client {
	var con Config
	con.Endpoint = config.Config.Trade.Endpoint
	con.WSEndpoint = config.Config.Trade.WSEndpoint

	con.ApiKey = apiKey
	con.SecretKey = secretKey
	con.Passphrase = passphrase
	con.TimeoutSecond = config.Config.Trade.TimeoutSecond
	con.IsPrint = config.Config.Trade.IsPrint
	con.I18n = config.Config.Trade.I18n

	client := NewClient(con)
	return client
}
