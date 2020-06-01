package tasks

import (
	"EX_okexquant/data"
	"EX_okexquant/db"
	"EX_okexquant/mylog"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRouter(r *gin.Engine) {
	/****************************** 永续合约 *********************************/
	r.POST("/api/swap/position/:instrument_id", GetSwapInstrumentPosition)
	r.POST("/api/swap/accounts/:instrument_id", GetSwapInstrumentAccount)
	r.POST("/api/swap/ledger/:instrument_id", GetSwapInstrumentLedger)
	r.POST("/api/swap/order", PostSwapOrder)
	r.POST("/api/swap/cancel_order/:instrument_id/:order_id", CancelSwapInstrumentOrder)
	r.POST("/api/swap/orders/:user_id/:instrument_id", GetSwapOrders)
	r.POST("/api/swap/fills/:instrument_id/:order_id", GetSwapFills)
}

/**
获取单个合约持仓信息
*/
func GetSwapInstrumentPosition(c *gin.Context) {
	out := data.CommonResp{}

	token := c.GetHeader("token")
	userID, err := db.ConvertTokenToUserID(token)
	instrumentID := c.Param("instrument_id")

	mylog.Logger.Info().Msgf("[Task Service] GetSwapInstrumentPosition request param: %s, %s", userID, instrumentID)

	if err != nil || instrumentID == "" {
		out.ErrorCode = data.EC_PARAMS_ERR
		out.ErrorMessage = data.ErrorCodeMessage(data.EC_PARAMS_ERR)
		c.JSON(http.StatusBadRequest, out)
		return
	}

	client, err := db.GetClientByUserID(userID)
	if err != nil {
		out.ErrorCode = data.EC_NETWORK_ERR
		out.ErrorMessage = err.Error()
		c.JSON(http.StatusBadRequest, out)
		return
	}

	list, err := client.GetSwapPositionByInstrument(instrumentID)
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
单个币种合约账户信息
*/
func GetSwapInstrumentAccount(c *gin.Context) {
	out := data.CommonResp{}

	token := c.GetHeader("token")
	userID, err := db.ConvertTokenToUserID(token)
	instrumentID := c.Param("instrumentID")

	mylog.Logger.Info().Msgf("[Task Service] GetSwapInstrumentAccount request param: %s, %s", userID, instrumentID)

	if err != nil || instrumentID == "" {
		out.ErrorCode = data.EC_PARAMS_ERR
		out.ErrorMessage = data.ErrorCodeMessage(data.EC_PARAMS_ERR)
		c.JSON(http.StatusBadRequest, out)
		return
	}

	client, err := db.GetClientByUserID(userID)
	if err != nil {
		out.ErrorCode = data.EC_NETWORK_ERR
		out.ErrorMessage = err.Error()
		c.JSON(http.StatusBadRequest, out)
		return
	}

	list, err := client.GetSwapAccount(instrumentID)
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
获取某个合约的用户配置
*/
func GetSwapInstrumentLeverage(c *gin.Context) {
	out := data.CommonResp{}

	token := c.GetHeader("token")
	userID, err := db.ConvertTokenToUserID(token)
	instrumentID := c.Param("instrumentID")

	mylog.Logger.Info().Msgf("[Task Service] GetSwapInstrumentLeverage request param: %s, %s", userID, instrumentID)

	if err != nil || instrumentID == "" {
		out.ErrorCode = data.EC_PARAMS_ERR
		out.ErrorMessage = data.ErrorCodeMessage(data.EC_PARAMS_ERR)
		c.JSON(http.StatusBadRequest, out)
		return
	}

	client, err := db.GetClientByUserID(userID)
	if err != nil {
		out.ErrorCode = data.EC_NETWORK_ERR
		out.ErrorMessage = err.Error()
		c.JSON(http.StatusBadRequest, out)
		return
	}

	list, err := client.GetSwapAccountsSettingsByInstrument(instrumentID)
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
设置某个合约的杠杆
*/
func SetSwapInstrumentLeverage(c *gin.Context) {
	out := data.CommonResp{}

	token := c.GetHeader("token")
	userID, err := db.ConvertTokenToUserID(token)
	instrumentID := c.Param("instrumentID")
	leverage := c.Param("leverage")
	side := c.Param("side")

	mylog.Logger.Info().Msgf("[Task Service] SetSwapInstrumentLeverage request param: %s, %s, %s, %s",
		userID, instrumentID, leverage, side)

	if err != nil || instrumentID == "" || leverage == "" || side == "" {
		out.ErrorCode = data.EC_PARAMS_ERR
		out.ErrorMessage = data.ErrorCodeMessage(data.EC_PARAMS_ERR)
		c.JSON(http.StatusBadRequest, out)
		return
	}

	client, err := db.GetClientByUserID(userID)
	if err != nil {
		out.ErrorCode = data.EC_NETWORK_ERR
		out.ErrorMessage = err.Error()
		c.JSON(http.StatusBadRequest, out)
		return
	}

	list, err := client.PostSwapAccountsLeverage(instrumentID, leverage, side)
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
func GetSwapInstrumentLedger(c *gin.Context) {
	out := data.CommonResp{}

	token := c.GetHeader("token")
	userID, err := db.ConvertTokenToUserID(token)
	instrumentID := c.Param("instrumentID")
	after := c.Param("after")
	before := c.Param("before")
	limit := c.Param("limit")
	oType := c.Param("type")

	mylog.Logger.Info().Msgf("[Task Service] GetSwapUnderlyingLedger request param: %s, %s, %s, %s, %s, %s",
		userID, instrumentID, after, before, limit, oType)

	if err != nil || instrumentID == "" {
		out.ErrorCode = data.EC_PARAMS_ERR
		out.ErrorMessage = data.ErrorCodeMessage(data.EC_PARAMS_ERR)
		c.JSON(http.StatusBadRequest, out)
		return
	}

	client, err := db.GetClientByUserID(userID)
	if err != nil {
		out.ErrorCode = data.EC_NETWORK_ERR
		out.ErrorMessage = err.Error()
		c.JSON(http.StatusBadRequest, out)
		return
	}

	optionalParams := make(map[string]string)
	optionalParams["after"] = after
	optionalParams["before"] = before
	optionalParams["limit"] = limit
	optionalParams["type"] = oType

	list, err := client.GetSwapAccountLedger(instrumentID, optionalParams)
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
获取所有订单列表
*/
func GetSwapOrders(c *gin.Context) {
	out := data.CommonResp{}

	token := c.GetHeader("token")
	userID, err := db.ConvertTokenToUserID(token)
	instrumentID := c.Param("instrument_id")
	state := c.Param("state")
	after := c.Param("after")
	before := c.Param("before")
	limit := c.Param("limit")

	if err != nil || instrumentID == "" || state == "" {
		out.ErrorCode = data.EC_PARAMS_ERR
		out.ErrorMessage = data.ErrorCodeMessage(data.EC_PARAMS_ERR)
		c.JSON(http.StatusBadRequest, out)
		return
	}

	mylog.Logger.Info().Msgf("[Task Service] GetSwapOrders request param: %s, %s, %s", userID, token, instrumentID)

	client, err := db.GetClientByUserID(userID)
	if err != nil {
		out.ErrorCode = data.EC_NETWORK_ERR
		out.ErrorMessage = err.Error()
		c.JSON(http.StatusBadRequest, out)
		return
	}

	optionParam := make(map[string]string)
	optionParam["after"] = after
	optionParam["before"] = before
	optionParam["limit"] = limit

	list, err := client.GetSwapOrderByInstrumentId(instrumentID, state, optionParam)
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
获取成交明细
*/
func GetSwapFills(c *gin.Context) {
	out := data.CommonResp{}

	token := c.GetHeader("token")
	userID, err := db.ConvertTokenToUserID(token)
	instrumentID := c.Param("instrument_id")
	orderID := c.Param("order_id")
	after := c.Param("after")
	before := c.Param("before")
	limit := c.Param("limit")

	if err != nil || instrumentID == "" || orderID == "" {
		out.ErrorCode = data.EC_PARAMS_ERR
		out.ErrorMessage = data.ErrorCodeMessage(data.EC_PARAMS_ERR)
		c.JSON(http.StatusBadRequest, out)
		return
	}

	mylog.Logger.Info().Msgf("[Task Service] GetSwapFills request param: %s, %s, %s", userID, instrumentID, orderID)

	client, err := db.GetClientByUserID(userID)
	if err != nil {
		out.ErrorCode = data.EC_NETWORK_ERR
		out.ErrorMessage = err.Error()
		c.JSON(http.StatusBadRequest, out)
		return
	}

	optionParam := make(map[string]string)
	optionParam["after"] = after
	optionParam["before"] = before
	optionParam["limit"] = limit

	list, err := client.GetSwapFills(instrumentID, orderID, optionParam)
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
