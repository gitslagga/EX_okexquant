package trade

import (
	"EX_okexquant/config"
	"EX_okexquant/mylog"
)

var okexClient *Client

func Init() {
	okexClient = newOKExClient()
	mylog.Logger.Error().Msgf("newOKExClient Init success")
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
