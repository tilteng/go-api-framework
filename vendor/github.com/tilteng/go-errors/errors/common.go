package errors

var ErrJSONSchemaValidationFailed = &ErrorClass{
	Name:        "ErrJSONSchemaValidationFailed",
	Code:        "ERR_ID_BAD_DATA",
	Status:      400,
	Message:     "Invalid data provided",
	Description: "Validation against json schema failed",
}

var ErrInternalServerError = &ErrorClass{
	Name:        "ErrInternalServerError",
	Code:        "ERR_ID_INTERNAL_SERVER_ERROR",
	Status:      500,
	Message:     "An unhandled exception has occurred",
	Description: "Something bad happened. Please try again.",
}

var ErrInternalError = &ErrorClass{
	Name:        "ErrInternalError",
	Code:        "ERR_ID_INTERNAL_ERROR",
	Status:      500,
	Message:     "An unknown error has occurred",
	Description: "Something bad happened. Please try again.",
}
