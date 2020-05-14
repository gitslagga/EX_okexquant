package main

import (
	"EX_okexquant/config"
	"EX_okexquant/data"
	"EX_okexquant/db"
	"EX_okexquant/mylog"
	"EX_okexquant/proxy"
	"EX_okexquant/tasks"
	"EX_okexquant/trade"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var configAddr = flag.String("config", "./config/config.toml", "base configuration files for server")

func main() {
	var err error

	flag.Parse()

	if *configAddr == "" {
		panic("Configuration file path is not set, server exit")
	}
	config.LoadConfig(*configAddr)

	data.Location, err = time.LoadLocation(config.Config.Server.Location)
	if err != nil {
		panic(fmt.Sprintf("load location failed, err=%v", err))
	}
	mylog.ConfigLoggers()

	trade.Init()
	proxy.Init()

	db.InitRedisCli()
	db.InitRedisData()
	defer db.CloseRedisCli()
	db.InitMysqlCli()
	defer db.CloseMysqlCli()

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	tasks.InitRouter(r)

	srv := &http.Server{
		Addr:    config.Config.Server.Address,
		Handler: r,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("listen err:%v\n", err)
		}
	}()
	fmt.Println("the server start succeed!!!")

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGTERM)
	sg := <-quit
	fmt.Printf("receive the signal:%v\n", sg)

	close(data.ShutdownChan)

	data.Wg.Wait()
	fmt.Println("wg return...")

	fmt.Println("main shutdown")
}
