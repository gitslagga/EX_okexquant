package trade

import (
	"EX_okexquant/config"
	"EX_okexquant/mylog"
	"github.com/okex/V3-Open-API-SDK/okex-go-sdk-api"
)

var okexClient *okex.Client

func Init() {
	okexClient = newOKExClient()
	mylog.Logger.Error().Msgf("newOKExClient Init success")
}

func newOKExClient() *okex.Client {
	var con okex.Config
	con.Endpoint = config.Config.Trade.Endpoint
	con.WSEndpoint = config.Config.Trade.WSEndpoint

	con.ApiKey = config.Config.Trade.ApiKey
	con.SecretKey = config.Config.Trade.SecretKey
	con.Passphrase = config.Config.Trade.Passphrase
	con.TimeoutSecond = config.Config.Trade.TimeoutSecond
	con.IsPrint = config.Config.Trade.IsPrint
	con.I18n = config.Config.Trade.I18n

	client := okex.NewClient(con)
	return client
}

func GetAccountCurrencies() (data string, err error) {

	data, err = "hello slagga", nil//okexClient.GetAccountCurrencies()
	if err != nil {
		mylog.Logger.Error().Msgf("okexClient GetAccountCurrencies failed: %v", err.Error())
	}

	return
}


