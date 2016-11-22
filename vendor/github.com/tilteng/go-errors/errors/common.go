package errors

var ErrJSONSchemaValidationFailed = &ErrorClass{
	Name:        "ErrJSONSchemaValidationFailed",
	Code:        "ERR_ID_BAD_DATA",
	Status:      400,
	Title:       "Invalid data provided",
	Description: "Validation against json schema failed",
}

var ErrInternalServerError = &ErrorClass{
	Name:        "ErrInternalServerError",
	Code:        "ERR_ID_INTERNAL_SERVER_ERROR",
	Status:      500,
	Title:       "An unhandled exception has occurred",
	Description: "Something bad happened. This is generally used for uncaught panics.",
}

var ErrInternalError = &ErrorClass{
	Name:        "ErrInternalError",
	Code:        "ERR_ID_INTERNAL_ERROR",
	Status:      500,
	Title:       "An unknown error has occurred",
	Description: "A random error for which we're not looking occurred. Most likely these are errors we're not necessarily expecting, but occurred.",
}
