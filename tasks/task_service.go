package tasks

import (
	"EX_okexquant/data"
	"EX_okexquant/trade"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRouter(r *gin.Engine) {
	r.POST("/quant/get_currency_list", GetCurrencyList)
}

func GetCurrencyList(c *gin.Context) {
	out := data.CommonResp{}

	list, _ := trade.GetAccountCurrencies()

	out.ErrorCode = data.EC_NONE.Code()
	out.ErrorMsg = data.EC_NONE.String()
	out.Data = list

	c.JSON(http.StatusOK, out)
	return
}
