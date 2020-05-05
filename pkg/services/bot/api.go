package bot

import (
	"time"

	"github.com/juju/errors"
	"go.uber.org/zap"
	tb "gopkg.in/tucnak/telebot.v2"

	"github.com/mgorozii/webhookdelia/pkg/config"

	"github.com/mgorozii/webhookdelia/pkg/models"
	"github.com/mgorozii/webhookdelia/pkg/services/store"
)

const (
	pollTimeout = 10 * time.Second
)

type Opts struct {
	Conf  config.Conf
	Log   *zap.Logger
	Store store.Service
}

type Service interface {
	Start()
	Send(record models.Record, message string) error
}

func newBot(conf config.Conf) (*tb.Bot, error) {
	return tb.NewBot(tb.Settings{
		URL:    conf.TelegramURL,
		Token:  conf.TelegramToken,
		Poller: &tb.LongPoller{Timeout: pollTimeout},
	})
}

func New(opts Opts) (Service, error) {
	bot, err := newBot(opts.Conf)
	if err != nil {
		return nil, errors.Trace(err)
	}

	s := &service{
		bot:   bot,
		store: opts.Store,
		log:   opts.Log,
		conf:  opts.Conf,
	}

	s.init()

	return s, nil
}
