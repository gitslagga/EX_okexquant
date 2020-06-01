package tasks

import (
	"EX_okexquant/db"
	"fmt"
	"time"
)

func InitFutures() {
	fmt.Println("[Tasks] swap init ...")

	StartSwapInstrumentsTask()

	fmt.Println("[Tasks] swap init success.")
}

func StartSwapInstrumentsTask() {
	db.SetSwapInstruments()

	go func() {
		timer := time.NewTicker(24 * time.Hour)
		for {
			select {
			case <-timer.C:
				db.SetSwapInstruments()
			}
		}
	}()

	fmt.Println("[Tasks] StartFixOrdersTask succeed.")
}
