package tasks

import (
	"EX_okexquant/config"
	"EX_okexquant/data"
	"EX_okexquant/mylog"
	"EX_okexquant/proxy"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

func FindAccountAssets(userID, size, currencyID, accountType string) bool {
	url := fmt.Sprintf("/assets/v1/api/findUserAssets?userId=%v&currencyId=%v&accountType=%v",
		userID, currencyID, accountType)

	mylog.Logger.Info().Msgf("[FindAccountAssets], url: %v", url)
	respBody, _, statusCode := proxy.Get(config.Config.Service.NotifyUrl, url, func(*http.Request) {})
	mylog.Logger.Info().Msgf("[FindAccountAssets], respBody: %v", string(respBody))
	if statusCode != 200 {
		mylog.Logger.Error().Msgf("[FindAccountAssets] failed, statusCode=%v", statusCode)
		return false
	}

	var resp data.ResponseFindAccount
	err := json.Unmarshal(respBody, &resp)
	if err != nil {
		mylog.Logger.Error().Msgf("[FindAccountAssets] Unmarshal error, err: %v", err)
		return false
	}

	if resp.RespCode != 1 {
		err = errors.New(resp.RespDesc)
		mylog.Logger.Error().Msgf("[FindAccountAssets] RespCode error, err: %v", err)
		return false
	}

	mylog.Logger.Info().Msgf("[FindAccountAssets] succeed: userID:%v, currencyID:%v, accountType:%v",
		userID, currencyID, accountType)

	num, err := strconv.ParseFloat(size, 64)
	if err != nil {
		mylog.Logger.Error().Msgf("[FindAccountAssets] ParseFloat error, err: %v", err)
		return false
	}

	amount := num * data.BTCContractVal
	if currencyID == "4" {
		amount = num * data.USDTContractVal
	}

	if amount > resp.RespData.Available {
		return false
	}

	return true
}
