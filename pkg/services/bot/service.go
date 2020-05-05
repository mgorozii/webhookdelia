package bot

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/juju/errors"
	"go.uber.org/zap"
	tb "gopkg.in/tucnak/telebot.v2"

	"github.com/mgorozii/webhookdelia/pkg/config"

	"github.com/mgorozii/webhookdelia/pkg/models"
	"github.com/mgorozii/webhookdelia/pkg/services/store"
)

type service struct {
	bot   *tb.Bot
	log   *zap.Logger
	store store.Service
	conf  config.Conf
}

func (s *service) Start() {
	s.bot.Start()
}

func (s *service) Send(record models.Record, message string) error {
	if _, err := s.bot.Send(record, message, tb.ModeHTML); err != nil {
		return errors.Trace(err)
	}

	return nil
}

func (s *service) init() {
	s.bot.Handle(tb.OnAddedToGroup, s.onStart)
	s.bot.Handle("/start", s.onStart)
	s.bot.Handle("/recreate", s.onCreateNew)
}

func (s *service) onStart(m *tb.Message) {
	// if connection already exists
	if record, err := s.store.GetByChatID(m.Chat.ID); err == nil {
		if err := s.Send(record, s.formatMessage(record)); err != nil {
			s.log.Error(
				"can't send message",
				zap.Any("record", record),
				zap.Error(err),
			)

			return
		}
	}

	// or create new
	s.onCreateNew(m)
}

func (s *service) onCreateNew(m *tb.Message) {
	record := models.Record{
		ChatID: m.Chat.ID,
		UUID:   uuid.New(),
	}

	if err := s.store.Put(record); err != nil {
		s.log.Error("can't save record", zap.Error(err))
		return
	}

	if err := s.Send(record, s.formatMessage(record)); err != nil {
		s.log.Error(
			"can't send message",
			zap.Any("record", record),
			zap.Error(err),
		)

		return
	}
}

func (s *service) formatMessage(record models.Record) string {
	url := fmt.Sprintf(
		"%s/send/%s",
		s.conf.PublicURL,
		record.UUID.String(),
	)

	return fmt.Sprintf(
		"webhook url: <a href=\"%s\">%s</a>"+
			"\n\n"+
			"example: <a href=\"%s?text=hello\">%s?text=hello</a>",
		url, url, url, url)
}
