package api_framework

import (
	"strings"

	"github.com/pborman/uuid"
)

type UUID struct {
	uuid.UUID
}

func GenUUID() *UUID {
	uuid := uuid.NewRandom()
	if uuid == nil {
		return nil
	} else {
		return &UUID{uuid}
	}
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

func UUIDFromString(s string) *UUID {
	uuid := uuid.Parse(s)
	if uuid == nil {
		return nil
	} else {
		return &UUID{uuid}
	}
}
