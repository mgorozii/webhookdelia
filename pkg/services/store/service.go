package store

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/juju/errors"
	"github.com/philippgille/gokv"

	"github.com/mgorozii/webhookdelia/pkg/models"
)

var (
	errNotFound = errors.New("record not found")
)

type service struct {
	internal gokv.Store
}

func (s *service) Close() error {
	return s.internal.Close()
}

func (s *service) Put(record models.Record) error {
	if err := s.internal.Delete(record.UUID.String()); err != nil {
		return errors.Trace(err)
	}

	if err := s.internal.Delete(fmt.Sprint(record.ChatID)); err != nil {
		return errors.Trace(err)
	}

	if err := s.internal.Set(record.UUID.String(), record); err != nil {
		return errors.Trace(err)
	}

	if err := s.internal.Set(fmt.Sprint(record.ChatID), record); err != nil {
		return errors.Trace(err)
	}

	return nil
}

func (s *service) GetByChatID(chatID int64) (models.Record, error) {
	record := models.Record{}

	found, err := s.internal.Get(fmt.Sprint(chatID), &record)
	if err != nil {
		return record, errors.Trace(err)
	} else if !found {
		return record, errors.Trace(errNotFound)
	}

	return record, nil
}

func (s *service) GetByUUID(uid uuid.UUID) (models.Record, error) {
	record := models.Record{}

	found, err := s.internal.Get(uid.String(), &record)
	if err != nil {
		return record, errors.Trace(err)
	} else if !found {
		return record, errors.Trace(errNotFound)
	}

	return record, nil
}
