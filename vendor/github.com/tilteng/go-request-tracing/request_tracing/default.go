package request_tracing

import (
	"strings"

	"github.com/pborman/uuid"
)

func defaultIDGenerator() string {
	return "REQ" + strings.ToUpper(
		strings.Replace(
			uuid.New(),
			"-",
			"",
			-1,
		),
	)
}

var defaultTraceManager = NewRequestTraceManager()
