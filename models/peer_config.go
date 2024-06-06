package models

import "github.com/google/uuid"

type PeerConfig struct {
	Id     uuid.UUID
	Config string
}
