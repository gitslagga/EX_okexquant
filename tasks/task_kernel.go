package tasks

import (
	"EX_okexquant/config"
	"EX_okexquant/data"
	"EX_okexquant/db"
	"EX_okexquant/mylog"
	"EX_okexquant/proxy"
	"EX_okexquant/trade"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

/**
下单
*/
func PostSwapOrder(c *gin.Context) {
	out := data.CommonResp{}
	orderParam := data.OrderParam{}

	if err := c.ShouldBindJSON(&orderParam); err != nil {
		mylog.Logger.Info().Msgf("[Task Service] PostSwapOrder request orderParam: %s", orderParam)
		out.ErrorCode = data.EC_PARAMS_ERR
		out.ErrorMessage = data.ErrorCodeMessage(data.EC_PARAMS_ERR)
		c.JSON(http.StatusBadRequest, out)
		return
	}

	token := c.GetHeader("token")
	userID, err := db.ConvertTokenToUserID(token)
	if err != nil {
		out.ErrorCode = data.EC_PARAMS_ERR
		out.ErrorMessage = data.ErrorCodeMessage(data.EC_PARAMS_ERR)
		c.JSON(http.StatusBadRequest, out)
		return
	}

	mylog.Logger.Info().Msgf("[Task Service] PostSwapOrder request param: %s, %s", userID, orderParam)

	//开多或者开空的时候
	if orderParam.Type == "1" || orderParam.Type == "2" {
		//判断用户余额
		valid := FindAccountAssets(userID, orderParam.Size, orderParam.InstrumentID)
		if !valid {
			mylog.Logger.Info().Msgf("[Task Service] PostSwapOrder FindAccountAssets valid: %v", valid)
			out.ErrorCode = data.EC_INTERNAL_ERR_DB
			out.ErrorMessage = errors.New("not enough swap assets").Error()
			c.JSON(http.StatusBadRequest, out)
			return
		}
	}

	client, err := db.GetClientByUserID(userID)
	if err != nil {
		out.ErrorCode = data.EC_NETWORK_ERR
		out.ErrorMessage = err.Error()
		c.JSON(http.StatusBadRequest, out)
		return
	}

	order := &trade.BasePlaceOrderInfo{}
	order.Size = orderParam.Size
	order.Type = orderParam.Type
	order.MatchPrice = orderParam.MatchPrice
	order.Price = orderParam.Price

	list, err := client.PostSwapOrder(orderParam.InstrumentID, order)
	if err != nil {
		out.ErrorCode = data.EC_NETWORK_ERR
		out.ErrorMessage = err.Error()
		c.JSON(http.StatusBadRequest, out)
		return
	}

	out.ErrorCode = data.EC_NONE.Code()
	out.ErrorMessage = data.EC_NONE.String()
	out.Data = list

	c.JSON(http.StatusOK, out)
	return
}

/**
撤单
*/
func CancelSwapInstrumentOrder(c *gin.Context) {
	out := data.CommonResp{}

	token := c.GetHeader("token")
	userID, err := db.ConvertTokenToUserID(token)
	instrumentID := c.Param("instrument_id")
	orderID := c.Param("order_id")

	if err != nil || instrumentID == "" || orderID == "" {
		out.ErrorCode = data.EC_PARAMS_ERR
		out.ErrorMessage = data.ErrorCodeMessage(data.EC_PARAMS_ERR)
		c.JSON(http.StatusBadRequest, out)
		return
	}

	mylog.Logger.Info().Msgf("[Task Service] CancelSwapInstrumentOrder request param: %s, %s, %s", userID, instrumentID, orderID)

	client, err := db.GetClientByUserID(userID)
	if err != nil {
		out.ErrorCode = data.EC_NETWORK_ERR
		out.ErrorMessage = err.Error()
		c.JSON(http.StatusBadRequest, out)
		return
	}

	list, err := client.PostSwapCancelOrder(instrumentID, orderID)
	if err != nil {
		out.ErrorCode = data.EC_NETWORK_ERR
		out.ErrorMessage = err.Error()
		c.JSON(http.StatusBadRequest, out)
		return
	}

	out.ErrorCode = data.EC_NONE.Code()
	out.ErrorMessage = data.EC_NONE.String()
	out.Data = list

	c.JSON(http.StatusOK, out)
	return
}

func FindAccountAssets(userID, size, instrumentID string) bool {
	instruments, err := db.GetSwapInstruments(instrumentID)
	if err != nil {
		return false
	}

	contractVal, err := strconv.ParseFloat(instruments["contractval"], 64)
	if err != nil {
		mylog.Logger.Error().Msgf("[FindAccountAssets] ParseFloat error, err: %v", err)
		return false
	}

	currencyID := "3"
	if instruments["quotecurrency"] == "USDT" {
		currencyID = "4"
	}

	url := fmt.Sprintf("/assets/v1/api/findUserAssets?userId=%v&currencyId=%v&accountType=%v",
		userID, currencyID, data.SwapContractType)

	mylog.Logger.Info().Msgf("[FindAccountAssets], url: %v", url)
	respBody, _, statusCode := proxy.Get(config.Config.Service.NotifyUrl, url, func(*http.Request) {})
	mylog.Logger.Info().Msgf("[FindAccountAssets], respBody: %v", string(respBody))
	if statusCode != 200 {
		mylog.Logger.Error().Msgf("[FindAccountAssets] failed, statusCode=%v", statusCode)
		return false
	}

	var resp data.ResponseFindAccount
	err = json.Unmarshal(respBody, &resp)
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
		userID, currencyID, data.SwapContractType)

	num, err := strconv.ParseFloat(size, 64)
	if err != nil {
		mylog.Logger.Error().Msgf("[FindAccountAssets] ParseFloat error, err: %v", err)
		return false
	}

	if (num / contractVal) >= resp.RespData.Available {
		return false
	}

	return true
}
