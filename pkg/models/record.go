package models

import (
	"fmt"

	"github.com/google/uuid"
)

type Record struct {
	ChatID int64
	UUID   uuid.UUID
}

func (r Record) Recipient() string {
	return fmt.Sprint(r.ChatID)
}
