package tasks

import (
	"EX_okexquant/db"
	"fmt"
	"time"
)

func InitFutures() {
	fmt.Println("[Tasks] futures init ...")

	StartFuturesInstrumentsTask()

	fmt.Println("[Tasks] futures init success.")
}

func StartFuturesInstrumentsTask() {
	db.GetFuturesInstruments()

	go func() {
		timer := time.NewTicker(24 * time.Hour)
		for {
			select {
			case <-timer.C:
				db.GetFuturesInstruments()
			}
		}
	}()

	fmt.Println("[Tasks] StartFixOrdersTask succeed.")
}
