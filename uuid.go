package goes

import (
	"github.com/satori/go.uuid"
)

// IDGen interface
//go:generate mockery -name IDGen
type IDGen interface {
	Generate() string
}

// UUIDV4 uuid version 4
type UUIDV4 struct{}

// NewUUIDV4 new uuid v1
func NewUUIDV4() *UUIDV4 {
	return &UUIDV4{}
}

// Generate uuid v1
func (u *UUIDV4) Generate() string {
	return uuid.NewV4().String()
}
