package ecode

// BadRequest new BadRequest error that is mapped to a 400 response.
func BadRequest(code int, message string) *Error {
	return New(code, 400, message)
}

// IsBadRequest determines if err is an error which indicates a BadRequest error.
// It supports wrapped errors.
func IsBadRequest(err error) bool {
	return HttpCode(err) == 400
}

// Unauthorized new Unauthorized error that is mapped to a 401 response.
func Unauthorized(code int, message string) *Error {
	return New(code, 401, message)
}

// IsUnauthorized determines if err is an error which indicates an Unauthorized error.
// It supports wrapped errors.
func IsUnauthorized(err error) bool {
	return HttpCode(err) == 401
}

// Forbidden new Forbidden error that is mapped to a 403 response.
func Forbidden(code int, message string) *Error {
	return New(code, 403, message)
}

// IsForbidden determines if err is an error which indicates a Forbidden error.
// It supports wrapped errors.
func IsForbidden(err error) bool {
	return HttpCode(err) == 403
}

func NotFoundWithMessage(code int, message string) *Error {
	return New(code, 404, message)
}

// IsNotFound determines if err is an error which indicates an NotFound error.
// It supports wrapped errors.
func IsNotFound(err error) bool {
	return HttpCode(err) == 404
}

// InternalServer new InternalServer error that is mapped to a 500 response.
func InternalServer(code int, message string) *Error {
	return New(code, 500, message)
}

// IsInternalServer determines if err is an error which indicates an Internal error.
// It supports wrapped errors.
func IsInternalServer(err error) bool {
	return HttpCode(err) == 500
}

// ClientClosed new ClientClosed error that is mapped to an HTTP 499 response.
func ClientClosed(code int, message string) *Error {
	return New(code, 499, message)
}

// IsClientClosed determines if err is an error which indicates a IsClientClosed error.
// It supports wrapped errors.
func IsClientClosed(err error) bool {
	return HttpCode(err) == 499
}
