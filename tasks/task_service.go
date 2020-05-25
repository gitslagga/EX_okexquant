package tasks

import (
	"EX_okexquant/data"
	"EX_okexquant/db"
	"EX_okexquant/mylog"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lithammer/shortuuid"
	"net/http"
	"strings"
)

func InitRouter(r *gin.Engine) {
	/****************************** 交割合约 *********************************/
	r.POST("/api/futures/position/:instrument_id", GetFuturesInstrumentPosition)
	r.POST("/api/futures/accounts/:underlying", GetFuturesUnderlyingAccount)
	r.POST("/api/futures/ledger/:underlying", GetFuturesUnderlyingLedger)
	r.POST("/api/futures/order", PostFuturesOrder)
	r.POST("/api/futures/cancel_order/:instrument_id/:order_id", CancelFuturesInstrumentOrder)
	r.POST("/api/futures/orders/:user_id/:instrument_id", GetFuturesOrders)
	r.POST("/api/futures/fills/:instrument_id/:order_id", GetFuturesFills)
}

/**
获取交割合约单个合约持仓信息
*/
func GetFuturesInstrumentPosition(c *gin.Context) {
	out := data.CommonResp{}
	instrumentID := c.Param("instrument_id")
	mylog.Logger.Info().Msgf("[Task Service] GetFuturesPosition request param: %s", instrumentID)

	if instrumentID == "" {
		out.ErrorCode = data.EC_PARAMS_ERR
		out.ErrorMessage = data.ErrorCodeMessage(data.EC_PARAMS_ERR)
		c.JSON(http.StatusBadRequest, out)
		return
	}

	list, err := db.GetFuturesInstrumentPosition(instrumentID)
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
交割合约下单
*/
func PostFuturesOrder(c *gin.Context) {
	out := data.CommonResp{}
	orderParam := data.OrderParam{}

	if err := c.ShouldBindJSON(&orderParam); err != nil {
		mylog.Logger.Info().Msgf("[Task Service] PostFuturesOrder request orderParam: %s", orderParam)
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

	optionParam := make(map[string]string)
	optionParam["client_oid"] = orderParam.ClientOID
	optionParam["order_type"] = orderParam.OrderType
	optionParam["match_price"] = orderParam.MatchPrice

	//开多或者开空的时候
	if orderParam.Type == "1" || orderParam.Type == "2" {
		//判断交易类型
		currencyID := "3"
		if strings.Contains(orderParam.InstrumentID, "USDT") {
			currencyID = "4"
		}

		//判断用户余额
		valid := FindAccountAssets(userID, orderParam.Size, currencyID, "3", "loginUserToken")
		if !valid {
			mylog.Logger.Info().Msgf("[Task Service] PostFuturesOrder FindAccountAssets valid: %v", valid)
			out.ErrorCode = data.EC_INTERNAL_ERR_DB
			out.ErrorMessage = errors.New("not enough futures assets").Error()
			c.JSON(http.StatusBadRequest, out)
			return
		}

		//生成通用唯一识别码
		optionParam["client_oid"] = shortuuid.New()
	}

	list, err := db.PostFuturesOrder(userID, orderParam.InstrumentID, orderParam.Type, orderParam.Price, orderParam.Size, optionParam)
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
交割合约撤单
*/
func CancelFuturesInstrumentOrder(c *gin.Context) {
	out := data.CommonResp{}

	instrumentID := c.Param("instrument_id")
	orderID := c.Param("order_id")
	mylog.Logger.Info().Msgf("[Task Service] CancelFuturesInstrumentOrder request param: %s, %s", instrumentID, orderID)

	if instrumentID == "" || orderID == "" {
		out.ErrorCode = data.EC_PARAMS_ERR
		out.ErrorMessage = data.ErrorCodeMessage(data.EC_PARAMS_ERR)
		c.JSON(http.StatusBadRequest, out)
		return
	}

	list, err := db.CancelFuturesInstrumentOrder(instrumentID, orderID)
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
交割合约获取订单列表
*/
func GetFuturesOrders(c *gin.Context) {
	out := data.CommonResp{}

	token := c.GetHeader("token")
	userID, err := db.ConvertTokenToUserID(token)
	if err != nil {
		out.ErrorCode = data.EC_PARAMS_ERR
		out.ErrorMessage = data.ErrorCodeMessage(data.EC_PARAMS_ERR)
		c.JSON(http.StatusBadRequest, out)
		return
	}

	instrumentID := c.Param("instrument_id")
	mylog.Logger.Info().Msgf("[Task Service] GetFuturesOrders request param: %s, %s", token, instrumentID)

	if userID == "" || instrumentID == "" {
		out.ErrorCode = data.EC_PARAMS_ERR
		out.ErrorMessage = data.ErrorCodeMessage(data.EC_PARAMS_ERR)
		c.JSON(http.StatusBadRequest, out)
		return
	}

	list, err := db.GetFuturesOrders(userID, instrumentID)
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
交割合约获取成交明细
*/
func GetFuturesFills(c *gin.Context) {
	out := data.CommonResp{}

	instrumentID := c.Param("instrument_id")
	orderID := c.Param("order_id")
	mylog.Logger.Info().Msgf("[Task Service] GetFuturesFills request param: %s, %s", instrumentID, orderID)

	if instrumentID == "" || orderID == "" {
		out.ErrorCode = data.EC_PARAMS_ERR
		out.ErrorMessage = data.ErrorCodeMessage(data.EC_PARAMS_ERR)
		c.JSON(http.StatusBadRequest, out)
		return
	}

	list, err := db.GetFuturesFills(instrumentID, orderID)
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
