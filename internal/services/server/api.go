package server

import (
	"go.uber.org/zap"

	"github.com/mgorozii/webhookdelia/internal/config"

	"github.com/mgorozii/webhookdelia/internal/services/bot"
	"github.com/mgorozii/webhookdelia/internal/services/store"
)

type Opts struct {
	Conf  config.Conf
	Bot   bot.Service
	Log   *zap.Logger
	Store store.Service
}

type Service interface {
	Run(addr ...string) error
}

func New(opts Opts) Service {
	s := &service{}
	s.init(opts)

	return s
}
