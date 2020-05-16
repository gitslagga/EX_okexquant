package tasks

import (
	"EX_okexquant/data"
	"EX_okexquant/db"
	"EX_okexquant/mylog"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRouter(r *gin.Engine) {
	/****************************** 交割合约 *********************************/
	r.POST("/api/futures/instruments/ticker", GetFuturesInstrumentsTicker)
	r.POST("/api/futures/position/:instrument_id", GetFuturesInstrumentPosition)
	r.POST("/api/futures/accounts/:underlying", GetFuturesUnderlyingAccount)
	r.POST("/api/futures/get/leverage/:underlying", GetFuturesUnderlyingLeverage)
	r.POST("/api/futures/set/leverage/:underlying/:leverage", SetFuturesUnderlyingLeverage)
	r.POST("/api/futures/ledger/:underlying", GetFuturesUnderlyingLedger)
}

/**
获取交割合约全部交易对
*/
func GetFuturesInstrumentsTicker(c *gin.Context) {
	out := data.CommonResp{}

	list, err := db.GetFuturesInstrumentsTicker()
	if err != nil {
		out.ErrorCode = data.EC_NETWORK_ERR
		out.ErrorMessage = data.ErrorCodeMessage(data.EC_NETWORK_ERR)
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
获取交割合约单个合约持仓信息
*/
func GetFuturesInstrumentPosition(c *gin.Context) {
	out := data.CommonResp{}
	instrumentID := c.Param("instrument_id")
	mylog.Logger.Info().Msgf("[Task Service] GetFuturesPosition request param: %s", instrumentID)

	if !db.ISExistsInstrumentsTicker(instrumentID) {
		out.ErrorCode = data.EC_PARAMS_ERR
		out.ErrorMessage = data.ErrorCodeMessage(data.EC_PARAMS_ERR)
		c.JSON(http.StatusBadRequest, out)
		return
	}

	list, err := db.GetFuturesInstrumentPosition(instrumentID)
	if err != nil {
		out.ErrorCode = data.EC_NETWORK_ERR
		out.ErrorMessage = data.ErrorCodeMessage(data.EC_NETWORK_ERR)
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
单个币种合约账户信息
*/
func GetFuturesUnderlyingAccount(c *gin.Context) {
	out := data.CommonResp{}

	underlying := c.Param("underlying")
	mylog.Logger.Info().Msgf("[Task Service] GetFuturesUnderlyingAccount request param: %s", underlying)

	if underlying == "" {
		out.ErrorCode = data.EC_PARAMS_ERR
		out.ErrorMessage = data.ErrorCodeMessage(data.EC_PARAMS_ERR)
		c.JSON(http.StatusBadRequest, out)
		return
	}

	list, err := db.GetFuturesUnderlyingAccount(underlying)
	if err != nil {
		out.ErrorCode = data.EC_NETWORK_ERR
		out.ErrorMessage = data.ErrorCodeMessage(data.EC_NETWORK_ERR)
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
获取交割合约杠杆倍数
*/
func GetFuturesUnderlyingLeverage(c *gin.Context) {
	out := data.CommonResp{}

	underlying := c.Param("underlying")
	mylog.Logger.Info().Msgf("[Task Service] GetFuturesUnderlyingLeverage request param: %s", underlying)

	if underlying == "" {
		out.ErrorCode = data.EC_PARAMS_ERR
		out.ErrorMessage = data.ErrorCodeMessage(data.EC_PARAMS_ERR)
		c.JSON(http.StatusBadRequest, out)
		return
	}

	list, err := db.GetFuturesUnderlyingLeverage(underlying)
	if err != nil {
		out.ErrorCode = data.EC_NETWORK_ERR
		out.ErrorMessage = data.ErrorCodeMessage(data.EC_NETWORK_ERR)
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
设置交割合约杠杆倍数
*/
func SetFuturesUnderlyingLeverage(c *gin.Context) {
	out := data.CommonResp{}

	underlying := c.Param("underlying")
	leverage := c.Param("leverage")
	mylog.Logger.Info().Msgf("[Task Service] SetFuturesUnderlyingLeverage request param: %s, %s", underlying, leverage)

	if underlying == "" || leverage == "" {
		out.ErrorCode = data.EC_PARAMS_ERR
		out.ErrorMessage = data.ErrorCodeMessage(data.EC_PARAMS_ERR)
		c.JSON(http.StatusBadRequest, out)
		return
	}

	list, err := db.SetFuturesUnderlyingLeverage(underlying, leverage)
	if err != nil {
		out.ErrorCode = data.EC_NETWORK_ERR
		out.ErrorMessage = data.ErrorCodeMessage(data.EC_NETWORK_ERR)
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
获取账单流水查询
*/
func GetFuturesUnderlyingLedger(c *gin.Context) {
	out := data.CommonResp{}

	underlying := c.Param("underlying")
	mylog.Logger.Info().Msgf("[Task Service] GetFuturesUnderlyingLedger request param: %s", underlying)

	if underlying == "" {
		out.ErrorCode = data.EC_PARAMS_ERR
		out.ErrorMessage = data.ErrorCodeMessage(data.EC_PARAMS_ERR)
		c.JSON(http.StatusBadRequest, out)
		return
	}

	list, err := db.GetFuturesUnderlyingLedger(underlying)
	if err != nil {
		out.ErrorCode = data.EC_NETWORK_ERR
		out.ErrorMessage = data.ErrorCodeMessage(data.EC_NETWORK_ERR)
		c.JSON(http.StatusBadRequest, out)
		return
	}

	out.ErrorCode = data.EC_NONE.Code()
	out.ErrorMessage = data.EC_NONE.String()
	out.Data = list

	c.JSON(http.StatusOK, out)
	return
}
