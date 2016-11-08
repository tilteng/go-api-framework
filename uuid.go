package framework

import (
	"strings"

	"github.com/pborman/uuid"
)

type UUID uuid.UUID

func GenUUID() UUID {
	return UUID(uuid.NewRandom())
}

func GenUUIDHex() string {
	return strings.ToUpper(
		strings.Replace(
			uuid.New(),
			"-",
			"",
			-1,
		),
	)
}
