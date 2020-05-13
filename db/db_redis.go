package db

import (
	"EX_okexquant/config"
	"EX_okexquant/mylog"
	"EX_okexquant/trade"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"time"
)

const (
	QuantCurrencyList               = "quant:currency:list"
)

var (
	redisPool *redis.Pool
)

func InitRedisCli() {
	address := config.Config.Redis.Address
	password := config.Config.Redis.Password
	maxActive := config.Config.Redis.MaxActive
	maxIdle := config.Config.Redis.MaxIdle
	idleMills := config.Config.Redis.IdleMills

	redisPool = newPool(address, password, maxIdle, maxActive, idleMills)
	_, err := redisPool.Dial()
	if err != nil {
		mylog.Logger.Fatal().Msgf("[InitRedis] dial redis failed, address=%v, password=%v", address, password)
	}

	fmt.Println("[Init] redis init succeed.")
}

func CloseRedisCli() {
	redisPool.Close()
}

func newPool(server, password string, maxidle int, maxactive int, idleMills int) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     maxidle,
		IdleTimeout: time.Duration(idleMills) * time.Millisecond,
		MaxActive:   maxactive,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				mylog.Logger.Error().Msgf("[Dial] Dial redis pool failed, err=%v", err)
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					mylog.Logger.Error().Msgf("[Dial] Auth redis cluster failed, err=%v", err)
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				mylog.Logger.Error().Msgf("[TestOnBorrow] Ping to redis cluster failed, err=%v", err)
			}
			return err
		},
	}
}

func InitRedisData() {

	//setcurrencylist
	boolen, err := CurrentCurrencyIsExist()
	if err != nil {
		panic(err)
	} else if boolen == false {
		currencyList, err := trade.GetAccountCurrencies()
		if err != nil {
			panic(err)
		}

		err = SetCurrentList(currencyList)
		if err != nil {
			panic(err)
		}
	}
}

func CurrentCurrencyIsExist() (bool, error) {
	redisConn := redisPool.Get()
	defer redisConn.Close()

	result, err := redis.Int(redisConn.Do("exists", QuantCurrencyList))
	if err != nil {
		mylog.Logger.Error().Msgf("redis get %v error, err:%v", QuantCurrencyList, err)
		return false, err
	}

	if result == 0  {
		return false, nil
	}

	return true, nil
}

func GetCurrencyList() (uint64, error) {
	redisConn := redisPool.Get()
	defer redisConn.Close()

	result, err := redis.Uint64(redisConn.Do("GET", QuantCurrencyList))
	if err != nil {
		mylog.Logger.Error().Msgf("redis GET %v error, err:%v", QuantCurrencyList, err)
		return 0, err
	}

	return result, err
}

func SetCurrentList(currencyList string) error {
	redisConn := redisPool.Get()
	defer redisConn.Close()

	_, err := redisConn.Do("SET", QuantCurrencyList, currencyList)
	if err != nil {
		mylog.Logger.Error().Msgf("redis SET %v error, err:%v", QuantCurrencyList, err)
		return err
	}

	return err
}
