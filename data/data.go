package data

import (
	"errors"
	"sync"
	"time"
)

var (
	Location     *time.Location
	Wg           sync.WaitGroup
	ShutdownChan = make(chan int)
)

type ErrorCode int

const (
	EC_NONE               ErrorCode = iota
	EC_PARAMS_ERR                   = 30110100
	EC_NETWORK_ERR                  = 30110101
	EC_INTERNAL_ERR                 = 30110102
	EC_INTERNAL_ERR_DB              = 30110103
	EC_INTERNAL_ERR_REDIS           = 30110104
)

func (c ErrorCode) Code() (r int) {
	r = int(c)
	return
}

func (c ErrorCode) Error() (r error) {
	r = errors.New(c.String())
	return
}

func (c ErrorCode) String() (r string) {
	switch c {
	case EC_NONE:
		r = "ok"
	case EC_NETWORK_ERR:
		r = "Network error"
	case EC_PARAMS_ERR:
		r = "Parameter error"
	case EC_INTERNAL_ERR:
		r = "Server error"
	case EC_INTERNAL_ERR_DB:
		r = "Server error"
	case EC_INTERNAL_ERR_REDIS:
		r = "Server error"

	default:
	}
	return
}

func ErrorCodeMessage(c ErrorCode) (r string) {
	return c.String()
}

type CommonResp struct {
	ErrorCode    int         `json:"error_code" form:"error_code"`
	ErrorMessage string      `json:"error_message" form:"error_message"`
	Data         interface{} `json:"data" form:"data"`
}

/*********************************** trading *************************************/
type OrderParam struct {
	ClientOID    string `json:"client_oid"`
	InstrumentID string `json:"instrument_id" binding:"required"`
	Type         string `json:"type" binding:"required"`
	Price        string `json:"price"`
	Size         string `json:"size" binding:"required"`
	OrderType    string `json:"order_type"`
	MatchPrice   string `json:"match_price"`
}

const (
	FuturesContractType int = 3
	SwapContractType    int = 5
)

type ResponseFindAccount struct {
	RespCode int         `json:"respCode"`
	RespDesc string      `json:"respDesc"`
	RespData accountInfo `json:"respData"`
}

type accountInfo struct {
	AssetsID     int     `json:"assetsId"`
	UserID       int     `json:"userId"`
	CurrencyID   int     `json:"currencyId"`
	CurrencyName string  `json:"currencyName"`
	Available    float64 `json:"available"`
	Freeze       float64 `json:"freeze"`
	Total        float64 `json:"total"`
	Give         float64 `json:"give"`
	CreateTime   string  `json:"createTime"`
	UpdateTime   string  `json:"updateTime"`
	Version      int     `json:"version"`
	IsLock       int     `json:"isLock"`
	UsdReal      int     `json:"usdReal"`
}
