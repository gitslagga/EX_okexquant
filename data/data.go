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

	EC_NO_ADDRESS                = 30100000 + 10
	EC_NO_BALANCE                = 30200000 + 10
	EC_BALANCE_UNCONFIRM         = 30200000 + 11
	EC_VALUE_NOT_STANDARD        = 30200000 + 12
	EC_BLOCK_ERR                 = 30200000 + 13
	EC_NO_TRANSACTION            = 30200000 + 14
	EC_ADDRESS_INVALID           = 30200000 + 15
	EC_STATUS_CHANGE             = 30200000 + 16
	EC_WITHDRAW_NO_FOUND         = 30200000 + 17
	EC_TOKEN_NO_FOUND            = 30200000 + 18
	EC_Main_Insufficient_balance = 30200000 + 19
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

	case EC_NO_ADDRESS:
		r = "No address available"
	case EC_NO_BALANCE:
		r = "Insufficient balance"
	case EC_BALANCE_UNCONFIRM:
		r = "Balance unconfirmed"
	case EC_VALUE_NOT_STANDARD:
		r = "Value not standard"
	case EC_BLOCK_ERR:
		r = "Block error"
	case EC_NO_TRANSACTION:
		r = "No transaction"
	case EC_ADDRESS_INVALID:
		r = "Invalid address"
	case EC_STATUS_CHANGE:
		r = "Status changed"
	case EC_WITHDRAW_NO_FOUND:
		r = "Withdraw id not found"
	case EC_TOKEN_NO_FOUND:
		r = "Token not found"
	case EC_Main_Insufficient_balance:
		r = "Insufficient funds"
	default:
	}
	return
}

func ErrorCodeMsg(c ErrorCode) (r string) {
	return c.String()
}

type CommonResp struct {
	ErrorCode int    `json:"err_code" form:"err_code"`
	ErrorMsg  string `json:"msg" form:"msg"`
	Data      string `json:"data" form:"data"`
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

