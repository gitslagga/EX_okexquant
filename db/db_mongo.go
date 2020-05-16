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
	err          error
	client       *mongo.Client
	collection   *mongo.Collection
	insertOneRes *mongo.InsertOneResult
	deleteRes    *mongo.DeleteResult
	updateRes    *mongo.UpdateResult
	cursor       *mongo.Cursor
	size         int64
)

func InitMongoCli() {
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

func GetFuturesInstrumentsTicker() (interface{}, error) {
	collection = client.Database("main_quantify").Collection("futures_instruments_ticker")
	cursor, err = collection.Find(getContext(), bson.D{})
	if err == mongo.ErrNoDocuments {
		list, err := trade.OKexClient.GetAccountCurrencies()
		if err != nil {
			mylog.Logger.Fatal().Msgf("[GetFuturesInstrumentsTicker] OKexClient GetAccountCurrencies failed, err=%v, list=%v", err, list)
			return nil, err
		}

		var dataArray []interface{}
		for i := 0; i < len(*list); i++ {
			dataArray = append(dataArray, (*list)[i])
		}

		_, _ = collection.InsertMany(getContext(), dataArray)
		return dataArray, nil
	}

	if err != nil {
		mylog.Logger.Fatal().Msgf("[GetFuturesInstrumentsTicker] collection Find failed, err=%v, cursor=%v", err, cursor)
		return nil, err
	}
	defer cursor.Close(context.Background())

	var record map[string]string
	var dataArray []map[string]string
	for cursor.Next(context.Background()) {
		_ = cursor.Decode(&record)
		dataArray = append(dataArray, record)
	}

	return dataArray, nil
}

func ISExistsInstrumentsTicker(instrumentID string) bool {
	collection = client.Database("main_quantify").Collection("futures_instruments_ticker")
	size, err = collection.CountDocuments(getContext(), bson.D{{"instrument_id", instrumentID}})

	if err != nil {
		mylog.Logger.Fatal().Msgf("[ISExistsInstrumentsTicker] collection CountDocuments failed, err=%v, collection=%v", err, collection)
		return false
	}

	if size <= 0 {
		return false
	}

	return true
}

func GetFuturesInstrumentPosition(instrumentID string) (interface{}, error) {
	position, err := trade.OKexClient.GetFuturesInstrumentPosition(instrumentID)
	if err != nil {
		mylog.Logger.Fatal().Msgf("[GetFuturesInstrumentsPosition] trade OKexClient failed, err=%v, position=%v", err, position)
		return nil, err
	}

	collection = client.Database("main_quantify").Collection("futures_instruments_position")
	size, err = collection.CountDocuments(getContext(), bson.D{})
	if err != nil {
		mylog.Logger.Fatal().Msgf("[InsertFuturesInstrumentsTicker] collection CountDocuments failed, err=%v, collection=%v", err, collection)
		return nil, err
	}

	if size <= 0 {
		_, _ = collection.InsertOne(getContext(), *position)
	}
	_, _ = collection.UpdateOne(getContext(), bson.D{{"instrument_id", instrumentID}}, *position)

	return *position, nil
}

func GetFuturesUnderlyingLeverage(underlying string) (interface{}, error) {
	leverage, err := trade.OKexClient.GetFuturesAccountsLeverage(underlying)
	if err != nil {
		mylog.Logger.Fatal().Msgf("[GetFuturesInstrumentsPosition] trade OKexClient failed, err=%v, leverage=%v", err, leverage)
		return nil, err
	}

	collection = client.Database("main_quantify").Collection("futures_instruments_leverage")
	size, err = collection.CountDocuments(getContext(), bson.D{})
	if size <= 0 {
		_, _ = collection.InsertOne(getContext(), leverage)
	}
	_, _ = collection.UpdateOne(getContext(), bson.D{{"underlying", underlying}}, leverage)

	return leverage, nil
}

func SetFuturesUnderlyingLeverage(underlying, leverage string) (interface{}, error) {
	resp, err := trade.OKexClient.PostFuturesAccountsLeverage(underlying, leverage, nil)
	if err != nil {
		mylog.Logger.Fatal().Msgf("[SetFuturesUnderlyingLeverage] trade OKexClient failed, err=%v, resp=%v", err, resp)
		return nil, err
	}

	collection = client.Database("main_quantify").Collection("futures_instruments_leverage")
	size, err = collection.CountDocuments(getContext(), bson.D{})
	if size <= 0 {
		_, _ = collection.InsertOne(getContext(), resp)
	}
	_, _ = collection.UpdateOne(getContext(), bson.D{{"underlying", underlying}}, resp)

	return resp, nil
}
