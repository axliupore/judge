package main

import (
	"fmt"
	"github.com/axliupore/judge/config"
	"github.com/axliupore/judge/handle"
	"github.com/axliupore/judge/pkg/log"
	"github.com/axliupore/judge/pkg/middleware"
	"github.com/axliupore/judge/pkg/nsq"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	log.InitLogger()

	r := gin.Default()

	r.Use(middleware.Cors())

	r.POST("/judge", handle.JudgeServer)
	r.POST("/exec", handle.ExecServer)

	go func() {
		err = r.Run(fmt.Sprintf(":%d", config.CoreConfig.Server.Port))
		if err != nil {
			log.Logger.Error("start server error", zap.Error(err))
		}
	}()

	consumer, err := nsq.NewConsumer(config.CoreConfig.Nsq.Topic, config.CoreConfig.Nsq.Channel)
	if err != nil {
		panic(err)
	}

	consumer.AddHandler(&handle.MessageHandler{})

	if err = consumer.ConnectToNSQLookupd(fmt.Sprintf("%s:%d", config.CoreConfig.Nsq.Address, config.CoreConfig.Nsq.Nsqlookupd)); err != nil {
		log.Logger.Errorf("Could not connect to nsqlookupd: %v", err)
	}

	<-consumer.StopChan
}
