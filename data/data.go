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

//mofify_account
type RequestModifyAccount struct {
	UserId  uint64  `json:"userId"`
	Address string  `json:"address"`
	Number  float64 `json:"number"`
	Chain   string  `json:"chain"`
	Coin    string  `json:"coin"`
	Txid    string  `json:"txid"`
}

type ResponseModifyAccount struct {
	RespCode int    `json:"respCode"`
	RespDesc string `json:"respDesc"`
	RespData string `json:"respData"`
}
