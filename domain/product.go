package domain

import (
	"errors"
	"time"
)

type client struct {
	ID          string
	CompanyName string
	CVR         string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (client *client) validate() error {
	if client.CompanyName == "" {
		return errors.New("Name cant be empty")
	}
	if len(client.CVR) != 8 {
		return errors.New("Has to be a valid CVR")
	}
	return nil
}
