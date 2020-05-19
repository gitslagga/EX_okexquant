package tasks

import (
	"fmt"
)

func InitFutures() {
	fmt.Println("[Tasks] futures init ...")

	//StartGetOrdersTask()

	fmt.Println("[Tasks] futures init success.")
}

//
//func StartGetOrdersTask() {
//	db.InsertFuturesInstrumentsTicker()
//
//	go func() {
//		timer := time.NewTicker(1 * time.Second)
//		for {
//			select {
//			case <-timer.C:
//				db.InsertFuturesInstrumentsTicker()
//			}
//		}
//	}()
//
//	fmt.Println("[Tasks] StartGetOrdersTask succeed.")
//}
