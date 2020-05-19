package tasks

import (
	"EX_okexquant/db"
	"fmt"
	"time"
)

func InitFutures() {
	fmt.Println("[Tasks] futures init ...")

	StartFixOrdersTask()

	fmt.Println("[Tasks] futures init success.")
}

//TODO fix order state task
func StartFixOrdersTask() {

	go func() {
		timer := time.NewTicker(2 * time.Second)
		for {
			select {
			case <-timer.C:
				db.FixFuturesInstrumentsOrders()
			}
		}
	}()

	fmt.Println("[Tasks] StartFixOrdersTask succeed.")
}
