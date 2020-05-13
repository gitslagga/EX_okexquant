package utils

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"errors"
)


/******float 精度在小数点8位******/

const (
	PRICE_FLOAT_FORMAT   = "%v.%08d"
	PRICE_DIV            = 100000000
	PRICE_BITLEN         = 8
	MAX_VALIED           = 18446744073709551615 //2*64
	MAX_REALPRICE_VALIED = 92233720368
)

func GetFloat64(uprice uint64) string {
	pricehigh := uprice / PRICE_DIV
	pricelow := uprice % PRICE_DIV

	sprice := fmt.Sprintf(PRICE_FLOAT_FORMAT, pricehigh, pricelow)

	return sprice
}

//a*10+e8
func GetUint(digits string) (c uint64, err error) {
	priceparts := strings.Split(strings.TrimSpace(digits), ".")
	if len(priceparts) == 1 {
		c, err = strconv.ParseUint(digits, 10, 64)
		if err != nil {
			return
		} else if c > MAX_REALPRICE_VALIED {
			err = errors.New("invalid value")
			return
		}

		c *= PRICE_DIV
		return
	} else if len(priceparts) == 2 {
		var uprice1 uint64
		var uprice2 uint64
		uprice1, err = strconv.ParseUint(priceparts[0], 10, 64)
		if err != nil {
			return
		} else if uprice1 > MAX_REALPRICE_VALIED {
			err = errors.New("invalid value")
			return
		}

		if len(priceparts[1]) <= 0 {
			c = uprice1 * PRICE_DIV
			return
		}

		if len(priceparts[1]) > PRICE_BITLEN {
			priceparts[1] = priceparts[1][:PRICE_BITLEN]
		}
		uprice2, err = strconv.ParseUint(priceparts[1], 10, 64)
		if err != nil {
			return
		}

		c = uprice1*PRICE_DIV + uprice2*uint64(math.Pow10(PRICE_BITLEN-len(priceparts[1])))
		return
	} else {
		err = errors.New("invalid value")
		return
	}

	return
}
