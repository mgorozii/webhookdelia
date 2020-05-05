package store

import (
	"github.com/google/uuid"
	"github.com/juju/errors"

	"github.com/mgorozii/webhookdelia/pkg/config"
	"github.com/mgorozii/webhookdelia/pkg/models"
)

type Service interface {
	Put(record models.Record) error
	GetByChatID(chatID int64) (models.Record, error)
	GetByUUID(uuid uuid.UUID) (models.Record, error)
	Close() error
}

type Opts struct {
	Conf config.Conf
}

func New(opts Opts) (Service, error) {
	internal, err := newInternalStore(opts.Conf)
	if err != nil {
		return &service{}, errors.Trace(err)
	}

	return &service{internal: internal}, nil
}
