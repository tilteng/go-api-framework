package errors

import (
	"strings"

	"github.com/pborman/uuid"
)

var defaultErrorIDGenerator = ErrIDGeneratorFn(func() string {
	return "ERR" + strings.ToUpper(
		strings.Replace(
			uuid.New(),
			"-",
			"",
			-1,
		),
	)
})

var defaultErrorManager = NewErrorManager()

func NewErrorClass(name string, code string, status int, title string) *ErrorClass {
	return defaultErrorManager.NewClass(name, code, status, title)
}
