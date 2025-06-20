package model

import (
	"github.com/oklog/ulid/v2"
)

type Auth struct {
	ID      ulid.ULID
	IsAdmin bool
}
