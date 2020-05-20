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

func GetFuturesInstrumentPosition(instrumentID string) (interface{}, error) {
	ctx := context.Background()
	position, err := trade.OKexClient.GetFuturesInstrumentPosition(instrumentID)
	if err != nil {
		mylog.Logger.Error().Msgf("[GetFuturesInstrumentsPosition] trade OKexClient failed, err=%v, position=%v", err, position)
		return nil, err
	}

	collection := client.Database("main_quantify").Collection("futures_instruments_position")
	size, err := collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		mylog.Logger.Error().Msgf("[InsertFuturesInstrumentsTicker] collection CountDocuments failed, err=%v", err)
		return nil, err
	}

	if size <= 0 {
		insertResult, err := collection.InsertOne(ctx, *position)
		if err != nil {
			mylog.Logger.Error().Msgf("[InsertFuturesInstrumentsTicker] collection InsertOne failed, err=%v, insertResult:%v", err, insertResult)
		}
	}

	updateResult, err := collection.UpdateOne(ctx, bson.D{{"instrument_id", instrumentID}}, bson.D{{"$set", *position}})
	if err != nil {
		mylog.Logger.Error().Msgf("[InsertFuturesInstrumentsTicker] collection UpdateOne failed, err=%v, updateResult:%v", err, updateResult)
	}

	return *position, nil
}

func GetFuturesUnderlyingAccount(underlying string) (interface{}, error) {
	ctx := context.Background()
	account, err := trade.OKexClient.GetFuturesAccountsByCurrency(underlying)
	if err != nil {
		mylog.Logger.Error().Msgf("[GetFuturesInstrumentsPosition] trade OKexClient failed, err=%v, account=%v", err, account)
		return nil, err
	}

	collection := client.Database("main_quantify").Collection("futures_underlying_account")
	size, err := collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		mylog.Logger.Error().Msgf("[GetFuturesInstrumentsPosition] collection CountDocuments failed, err=%v", err)
		return nil, err
	}

	if size <= 0 {
		insertResult, err := collection.InsertOne(ctx, account)
		if err != nil {
			mylog.Logger.Error().Msgf("[GetFuturesInstrumentsPosition] collection InsertOne failed, err=%v, insertResult:%v", err, insertResult)
		}
	}

	updateResult, err := collection.UpdateOne(ctx, bson.D{{"underlying", underlying}}, bson.D{{"$set", account}})
	if err != nil {
		mylog.Logger.Error().Msgf("[GetFuturesInstrumentsPosition] collection UpdateOne failed, err=%v, updateResult:%v", err, updateResult)
	}

	return account, nil
}

func GetFuturesUnderlyingLedger(underlying string) (interface{}, error) {
	ctx := context.Background()
	optionalParams := map[string]string{}
	optionalParams["limit"] = "100"

	ledger, err := trade.OKexClient.GetFuturesAccountsLedgerByCurrency(underlying, optionalParams)
	if err != nil {
		mylog.Logger.Error().Msgf("[GetFuturesUnderlyingLedger] trade OKexClient failed, err=%v, ledger=%v", err, ledger)
		return nil, err
	}

	var record map[string]interface{}
	collection := client.Database("main_quantify").Collection("futures_underlying_ledger")
	for _, v := range ledger {
		err := collection.FindOne(ctx, bson.D{
			{"ledger_id", v["ledger_id"]},
		}).Decode(&record)

		if err == mongo.ErrNoDocuments || len(record) <= 0 {
			insertResult, err := collection.InsertOne(ctx, v)
			if err != nil {
				mylog.Logger.Error().Msgf("[GetFuturesUnderlyingLedger] collection InsertOne failed, err=%v, insertResult:%v", err, insertResult)
			}
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
	collection := client.Database("main_quantify").Collection("futures_instruments_users")
	insertResult, err := collection.InsertOne(getContext(), *resp)
	if err != nil {
		mylog.Logger.Error().Msgf("[PostFuturesOrder] collection InsertOne failed, err=%v, insertResult:%v", err, insertResult)
	}

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
	ctx := context.Background()
	collection := client.Database("main_quantify").Collection("futures_instruments_users")
	cursor, err := collection.Find(ctx, bson.D{
		{"user_id", userID},
		{"instrument_id", instrumentID},
	}, options.Find().SetSort(bson.M{"_id": -1}), options.Find().SetLimit(100))
	if err != nil {
		mylog.Logger.Error().Msgf("[GetFuturesOrders] collection Find failed, err=%v, cursor=%v", err, cursor)
		return nil, err
	}

	defer cursor.Close(ctx)

	var recordArray []map[string]interface{}
	collection = client.Database("main_quantify").Collection("futures_instruments_orders")

	for cursor.Next(ctx) {
		orderID := cursor.Current.Lookup("order_id").StringValue()

		size, err := collection.CountDocuments(ctx, bson.D{
			{"instrument_id", instrumentID},
			{"order_id", orderID},
		})
		if err != nil {
			mylog.Logger.Error().Msgf("[GetFuturesOrders] collection CountDocuments failed, err=%v, size=%v", err, size)
		}

		if size <= 0 {
			insertFuturesInstrumentsOrder(ctx, instrumentID, orderID)
		}

		var order map[string]interface{}
		err = collection.FindOne(ctx, bson.D{
			{"instrument_id", instrumentID},
			{"order_id", orderID},
		}).Decode(&order)

		if err != nil {
			mylog.Logger.Error().Msgf("[GetFuturesOrders] collection FindOne failed, err=%v, order=%v", err, order)
		} else {
			recordArray = append(recordArray, order)
		}
	}

	return recordArray, nil
}

func GetFuturesFills(instrumentID, orderID string) (interface{}, error) {
	ctx := context.Background()
	collection := client.Database("main_quantify").Collection("futures_instruments_fills")
	size, err := collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		mylog.Logger.Error().Msgf("[GetFuturesFills] collection CountDocuments failed, err=%v, size=%v", err, size)
		return nil, err
	}
	if size <= 0 {
		insertFuturesInstrumentsFills(ctx, instrumentID, orderID)
	}

	cursor, err := collection.Find(ctx, bson.D{
		{"instrument_id", instrumentID},
		{"order_id", orderID},
	}, options.Find().SetSort(bson.M{"_id": -1}), options.Find().SetLimit(100))

	if err != nil {
		mylog.Logger.Error().Msgf("[GetFuturesFills] collection Find failed, err=%v, cursor=%v", err, cursor)
		return nil, err
	}

	defer cursor.Close(ctx)
	var recordArray []map[string]interface{}
	for cursor.Next(ctx) {
		var record map[string]interface{}
		err = cursor.Decode(&record)

		if err != nil {
			mylog.Logger.Error().Msgf("[GetFuturesFills] cursor Decode failed, err=%v, record=%v", err, record)
		} else {
			recordArray = append(recordArray, record)
		}
	}

	return recordArray, nil
}

func insertFuturesInstrumentsOrder(ctx context.Context, instrumentID, orderID string) {
	order, err := trade.OKexClient.GetFuturesOrder(instrumentID, orderID)
	if err != nil {
		mylog.Logger.Error().Msgf("[insertFuturesInstrumentsOrder] trade OKexClient failed, err:%v, order:%v", err, order)
	}

	if len(order) > 0 {
		collection := client.Database("main_quantify").Collection("futures_instruments_orders")
		insertResult, err := collection.InsertOne(ctx, order)
		if err != nil {
			mylog.Logger.Error().Msgf("[insertFuturesInstrumentsOrder] collection InsertOne failed, err:%v, insertResult:%v", err, insertResult)
		}
	}
}

func insertFuturesInstrumentsFills(ctx context.Context, instrumentID, orderID string) {
	optionalParams := map[string]string{}
	optionalParams["limit"] = "100"

	fills, err := trade.OKexClient.GetFuturesFills(instrumentID, orderID, optionalParams)
	if err != nil {
		mylog.Logger.Error().Msgf("[insertFuturesInstrumentsFills] trade OKexClient failed, err:%v, fills:%v", err, fills)
		return
	}

	var data []interface{}
	for _, v := range fills {
		data = append(data, v)
	}

	if len(data) > 0 {
		collection := client.Database("main_quantify").Collection("futures_instruments_fills")
		insertManyResult, err := collection.InsertMany(ctx, data)
		if err != nil {
			mylog.Logger.Error().Msgf("[insertFuturesInstrumentsFills] collection InsertMany failed, err:%v, insertManyResult:%v", err, insertManyResult)
		}
	}
}

func FixFuturesInstrumentsOrders() {
	ctx := context.Background()

	//对未成交，部分成交，下单中，撤单中的订单进行修正
	collection := client.Database("main_quantify").Collection("futures_instruments_orders")
	cursor, err := collection.Find(ctx, bson.D{{
		"state", bson.D{{"$in", bson.A{"0", "1", "3", "4"}}},
	}}, options.Find().SetSort(bson.M{"_id": -1}), options.Find().SetLimit(100))
	if err != nil {
		mylog.Logger.Error().Msgf("[FixFuturesInstrumentsOrders] collection Find failed, err=%v, cursor=%v", err, cursor)
	}

	defer cursor.Close(ctx)

	var record map[string]string
	for cursor.Next(ctx) {
		err = cursor.Decode(&record)

		if err != nil {
			mylog.Logger.Error().Msgf("[FixFuturesInstrumentsOrders] cursor Decode failed, err=%v, record=%v", err, record)
		} else {
			realOrder, err := trade.OKexClient.GetFuturesOrder(record["instrument_id"], record["order_id"])
			if err != nil {
				mylog.Logger.Error().Msgf("[FixFuturesInstrumentsOrders] trade OKexClient failed, err=%v, realOrder=%v", err, realOrder)
			}

			updateResult, err := collection.UpdateOne(ctx, bson.D{
				{"instrument_id", record["instrument_id"]},
				{"order_id", record["order_id"]},
			}, bson.D{{"$set", realOrder}})
			if err != nil {
				mylog.Logger.Error().Msgf("[FixFuturesInstrumentsOrders] collection UpdateOne failed, err=%v, updateResult=%v", err, updateResult)
			}
		}
	}
}
