package db

import (
	"EX_okexquant/config"
	"EX_okexquant/mylog"
	"EX_okexquant/trade"
	"context"
	"errors"
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

func GetFuturesInstrumentPosition(instrumentID string) (interface{}, error) {
	position, err := trade.OKexClient.GetFuturesInstrumentPosition(instrumentID)
	if err != nil {
		mylog.Logger.Error().Msgf("[GetFuturesInstrumentsPosition] trade OKexClient failed, err=%v, position=%v", err, position)
		return nil, err
	}

	collection = client.Database("main_quantify").Collection("futures_instruments_position")
	size, err = collection.CountDocuments(getContext(), bson.D{})
	if err != nil {
		mylog.Logger.Error().Msgf("[InsertFuturesInstrumentsTicker] collection CountDocuments failed, err=%v, collection=%v", err, collection)
		return nil, err
	}

	if size <= 0 {
		_, _ = collection.InsertOne(getContext(), *position)
	}
	_, _ = collection.UpdateOne(getContext(), bson.D{{"instrument_id", instrumentID}}, *position)

	return *position, nil
}

func GetFuturesUnderlyingAccount(underlying string) (interface{}, error) {
	account, err := trade.OKexClient.GetFuturesAccountsByCurrency(underlying)
	if err != nil {
		mylog.Logger.Error().Msgf("[GetFuturesInstrumentsPosition] trade OKexClient failed, err=%v, account=%v", err, account)
		return nil, err
	}

	collection = client.Database("main_quantify").Collection("futures_underlying_account")
	size, err = collection.CountDocuments(getContext(), bson.D{})
	if err != nil {
		mylog.Logger.Error().Msgf("[GetFuturesInstrumentsPosition] collection CountDocuments failed, err=%v, collection=%v", err, collection)
		return nil, err
	}

	if size <= 0 {
		_, _ = collection.InsertOne(getContext(), account)
	}
	_, _ = collection.UpdateOne(getContext(), bson.D{{"underlying", underlying}}, account)

	return account, nil
}

func GetFuturesUnderlyingLedger(underlying string) (interface{}, error) {
	optionalParams := map[string]string{}
	optionalParams["limit"] = "100"

	ledger, err := trade.OKexClient.GetFuturesAccountsLedgerByCurrency(underlying, optionalParams)
	if err != nil {
		mylog.Logger.Error().Msgf("[GetFuturesUnderlyingLedger] trade OKexClient failed, err=%v, ledger=%v", err, ledger)
		return nil, err
	}

	var record map[string]interface{}
	collection = client.Database("main_quantify").Collection("futures_underlying_ledger")
	for _, v := range ledger {
		err := collection.FindOne(getContext(), bson.D{
			{"ledger_id", v["ledger_id"]},
		}).Decode(&record)

		if err == mongo.ErrNoDocuments {
			_, _ = collection.InsertOne(getContext(), v)
		}
	}

	return ledger, nil
}

func PostFuturesOrder(userID, instrumentID, oType, price, size string, optionalParams map[string]string) (interface{}, error) {
	resp, err := trade.OKexClient.PostFuturesOrder(instrumentID, oType, price, size, optionalParams)
	if err != nil {
		mylog.Logger.Error().Msgf("[PostFuturesOrder] trade OKexClient failed, err=%v, order=%v", err, resp)
		return nil, err
	}

	if (*resp)["result"] != true {
		err = errors.New((*resp)["error_message"].(string))
		mylog.Logger.Error().Msgf("[PostFuturesOrder] trade OKexClient failed, err=%v, order=%v", err, resp)
		return nil, err
	}

	(*resp)["user_id"] = userID
	(*resp)["instrument_id"] = instrumentID
	collection = client.Database("main_quantify").Collection("futures_instruments_users")
	_, _ = collection.InsertOne(getContext(), *resp)

	return *resp, nil
}

func CancelFuturesInstrumentOrder(instrumentID, orderID string) (interface{}, error) {
	resp, err := trade.OKexClient.CancelFuturesInstrumentOrder(instrumentID, orderID)
	if err != nil {
		mylog.Logger.Error().Msgf("[CancelFuturesInstrumentOrder] trade OKexClient failed, err=%v, order=%v", err, resp)
		return nil, err
	}

	if resp["result"] != true {
		err = errors.New(resp["error_message"].(string))
		mylog.Logger.Error().Msgf("[PostFuturesOrder] trade OKexClient failed, err=%v, order=%v", err, resp)
		return nil, err
	}

	return resp, nil
}

func GetFuturesOrders(userID, instrumentID string) (interface{}, error) {

	collection = client.Database("main_quantify").Collection("futures_instruments_users")
	cursor, err = collection.Find(getContext(), bson.D{
		{"user_id", userID},
	})
	if err != nil {
		mylog.Logger.Error().Msgf("[GetFuturesOrders] collection Find failed, err=%v, cursor=%v", err, cursor)
		return nil, err
	}

	defer cursor.Close(context.Background())

	var record map[string]string
	var recordArray []map[string]string
	collection = client.Database("main_quantify").Collection("futures_instruments_orders")
	for cursor.Next(context.Background()) {
		_ = cursor.Decode(&record)

		size, _ = collection.CountDocuments(getContext(), bson.D{
			{"instrument_id", instrumentID},
			{"order_id", record["order_id"]},
		})
		if size <= 0 {
			insertFuturesInstrumentsOrder(instrumentID, record["order_id"])
		}

		var order map[string]string
		_ = collection.FindOne(getContext(), bson.D{
			{"instrument_id", instrumentID},
			{"order_id", record["order_id"]},
		}).Decode(&order)

		recordArray = append(recordArray, order)
	}

	return recordArray, nil
}

func GetFuturesFills(instrumentID, orderID string) (interface{}, error) {
	size, err = collection.CountDocuments(getContext(), bson.D{})
	if err != nil {
		mylog.Logger.Error().Msgf("[GetFuturesFills] collection CountDocuments failed, err=%v, size=%v", err, size)
		return nil, err
	}
	if size <= 0 {
		insertFuturesInstrumentsFills(instrumentID, orderID)
	}

	collection = client.Database("main_quantify").Collection("futures_instruments_fills")
	cursor, err = collection.Find(getContext(), bson.D{
		{"instrument_id", instrumentID},
		{"order_id", orderID},
	})

	if err != nil {
		mylog.Logger.Error().Msgf("[GetFuturesFills] collection Find failed, err=%v, cursor=%v", err, cursor)
		return nil, err
	}

	defer cursor.Close(context.Background())
	var record map[string]string
	var recordArray []map[string]string
	for cursor.Next(context.Background()) {
		_ = cursor.Decode(&record)
		recordArray = append(recordArray, record)
	}

	return recordArray, nil
}

func insertFuturesInstrumentsOrder(instrumentID, orderID string) {
	order, err := trade.OKexClient.GetFuturesOrder(instrumentID, orderID)
	if err != nil {
		mylog.Logger.Error().Msgf("insertFuturesInstrumentsOrder error! err:%v", err)
	}

	collection = client.Database("main_quantify").Collection("futures_instruments_orders")
	_, _ = collection.InsertOne(getContext(), order)
}

func insertFuturesInstrumentsFills(instrumentID, orderID string) {
	fills, err := trade.OKexClient.GetFuturesFills(instrumentID, orderID, nil)
	if err != nil {
		mylog.Logger.Error().Msgf("insertFuturesInstrumentsFills error! err:%v", err)
		return
	}

	var data []interface{}
	for _, v := range fills {
		data = append(data, v)
	}

	collection = client.Database("main_quantify").Collection("futures_instruments_orders")
	_, _ = collection.InsertMany(getContext(), data)
}

func FixFuturesInstrumentsOrders() {
	//对未成交，部分成交，下单中，撤单中的订单进行修正
	collection = client.Database("main_quantify").Collection("futures_instruments_orders")
	cursor, err = collection.Find(getContext(), bson.D{{
		"state", bson.D{{"$in", bson.A{"0", "1", "3", "4"}}},
	}})
	if err != nil {
		mylog.Logger.Error().Msgf("[FixFuturesInstrumentsOrders] collection Find failed, err=%v, cursor=%v", err, cursor)
	}

	defer cursor.Close(context.Background())
	var order map[string]string
	for cursor.Next(context.Background()) {
		_ = cursor.Decode(&order)

		realOrder, _ := trade.OKexClient.GetFuturesOrder(order["instrument_id"], order["order_id"])
		_, _ = collection.UpdateOne(getContext(), bson.D{
			{"instrument_id", order["instrument_id"]},
			{"order_id", order["order_id"]},
		}, realOrder)
	}
}
