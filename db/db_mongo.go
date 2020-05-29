package db

import (
	"EX_okexquant/config"
	"EX_okexquant/mylog"
	"EX_okexquant/trade"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var (
	client *mongo.Client
)

func InitMongoCli() {
	var err error
	uri := config.Config.Mongo.ApplyURI
	localThreshold := config.Config.Mongo.LocalThreshold
	maxConnIdleTime := config.Config.Mongo.MaxConnIdleTime
	maxPoolSize := config.Config.Mongo.MaxPoolSize

	opt := options.Client().ApplyURI(uri)
	opt.SetLocalThreshold(time.Duration(localThreshold) * time.Second)   //只使用与mongo操作耗时小于3秒的
	opt.SetMaxConnIdleTime(time.Duration(maxConnIdleTime) * time.Second) //指定连接可以保持空闲的最大毫秒数
	opt.SetMaxPoolSize(maxPoolSize)                                      //使用最大的连接数

	client, err = mongo.Connect(getContext(), opt)
	if err != nil {
		mylog.Logger.Fatal().Msgf("[InitMongoCli] mongo connection failed, err=%v, client=%v", err, client)
	}

	fmt.Println("[InitMongo] mongo succeed.")
}

func CloseMongoCli() {
	client.Disconnect(getContext())
}

func getContext() context.Context {
	timeout := config.Config.Mongo.Timeout
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)

	return ctx
}

/**
根据用户ID获取OKex客户端
*/
func GetClientByUserID(userID string) (*trade.Client, error) {
	ctx := context.Background()
	var userKeys map[string]string

	collection := client.Database("main_quantify").Collection("futures_user_keys")
	errUser := collection.FindOne(ctx, bson.D{{"user_id", userID}}).Decode(&userKeys)
	if errUser != nil {
		if errUser == mongo.ErrNoDocuments {
			err := collection.FindOne(ctx, bson.D{{"status", "0"}}).Decode(&userKeys)
			if err != nil {
				mylog.Logger.Error().Msgf("[GetClientByUserID] collection FindOne failed, err=%v", err)
				return nil, err
			}

			updateResult, err := collection.UpdateOne(ctx, bson.D{{"_id", userKeys["_id"]}}, bson.D{{
				"$set", bson.D{{"status", "1"}, {"user_id", userID}},
			}})
			if err != nil {
				mylog.Logger.Error().Msgf("[GetClientByUserID] collection UpdateOne failed, updateResult=%v, err=%v", updateResult, err)
				return nil, err
			}
		} else {
			mylog.Logger.Error().Msgf("[GetClientByUserID] collection FindOne failed, err=%v", errUser)
			return nil, errUser
		}
	}

	client := trade.NewClientByParam(userKeys["api_key"], userKeys["secret_key"], userKeys["passphrase"])
	return client, nil
}

/**
获取并设置合约信息
*/
func GetFuturesInstruments() {
	instruments, err := trade.OKexClient.GetFuturesInstruments()
	if err != nil {
		mylog.Logger.Error().Msgf("[GetFuturesInstruments] trade OKexClient failed, err:%v, instruments:%v", err, instruments)
		return
	}

	var data []interface{}
	for _, v := range instruments {
		data = append(data, v)
	}

	if len(data) > 0 {
		ctx := context.Background()
		collection := client.Database("main_quantify").Collection("futures_instruments")
		err = collection.Drop(ctx)
		if err != nil {
			mylog.Logger.Error().Msgf("[GetFuturesInstruments] collection Drop failed, err:%v", err)
		}

		insertManyResult, err := collection.InsertMany(ctx, data)
		if err != nil {
			mylog.Logger.Error().Msgf("[GetFuturesInstruments] collection InsertMany failed, err:%v, insertManyResult:%v", err, insertManyResult)
			return
		}
	}
}
