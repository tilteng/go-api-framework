package errors

// JSON Schema validation failed
var ErrJSONSchemaValidationFailed = NewErrorClass(
	"ErrJSONSchemaValidationFailed",
	"ERR_ID_BAD_DATA",
	400,
	"Invalid data provided",
)

// This is generally used for uncaught panics
var ErrInternalServerError = NewErrorClass(
	"ErrInternalServerError",
	"ERR_ID_INTERNAL_SERVER_ERROR",
	500,
	"An unhandled exception has occurred",
)

// A random error for which we're not looking occurred. Most likely these
// are errors we're not necessarily expecting, but occurred.",
var ErrInternalError = NewErrorClass(
	"ErrInternalError",
	"ERR_ID_INTERNAL_ERROR",
	500,
	"An unknown error has occurred",
)
