package errors

import (
	"strings"

	"github.com/pborman/uuid"
)

type ErrIDGenerator interface {
	GenErrID() string
}

type config struct {
	ErrIDGenerator ErrIDGenerator
}

type defaultErrIDGenerator struct{}

func (defaultErrIDGenerator) GenErrID() string {
	return "ERR" + strings.ToUpper(
		strings.Replace(
			uuid.New(),
			"-",
			"",
			-1,
		),
	)
}

var Config *config

func init() {
	Config = &config{
		ErrIDGenerator: &defaultErrIDGenerator{},
	}
}
