package tasks

import (
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	//r.POST("/quant/get_currency_list", GetCurrencyList)
}

//func GetCurrencyList(c *gin.Context) {
//	out := data.CommonResp{}
//
//	list, err := trade.GetAccountCurrencies()
//	if err != nil {
//		out.ErrorCode = data.EC_NETWORK_ERR
//		out.ErrorMessage = data.ErrorCodeMessage(data.EC_NETWORK_ERR)
//		c.JSON(http.StatusBadRequest, out)
//		return
//	}
//
//	out.ErrorCode = data.EC_NONE.Code()
//	out.ErrorMessage = data.EC_NONE.String()
//	out.Data = list
//
//	c.JSON(http.StatusOK, out)
//	return
//}
