package main

import (
	"go.uber.org/zap"

	"github.com/mgorozii/webhookdelia/pkg/config"
	"github.com/mgorozii/webhookdelia/pkg/services/bot"
	"github.com/mgorozii/webhookdelia/pkg/services/server"
	"github.com/mgorozii/webhookdelia/pkg/services/store"
)

func main() {
	conf, err := config.New()
	if err != nil {
		panic(err)
	}

	log, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	storeService, err := store.New(store.Opts{
		Conf: conf,
	})
	if storeService != nil {
		defer func() {
			if err := storeService.Close(); err != nil {
				log.Error("can't close store", zap.Error(err))
			}
		}()
	}

	if err != nil {
		log.Error("can't construct store", zap.Error(err))
		return
	}

	botService, err := bot.New(bot.Opts{
		Conf:  conf,
		Log:   log,
		Store: storeService,
	})
	if err != nil {
		log.Error("can't construct bot", zap.Error(err))
		return
	}

	srv := server.New(server.Opts{
		Conf:  conf,
		Bot:   botService,
		Log:   log,
		Store: storeService,
	})

	go botService.Start()

	if err := srv.Run(":" + conf.Port); err != nil {
		log.Error("server return error", zap.Error(err))
	}
}
