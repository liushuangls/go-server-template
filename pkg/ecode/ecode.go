package ecode

var (
	InternalServerErr = InternalServer(-1, "Internal Server Error")

	InvalidParams  = NewInvalidParamsErr("Invalid Params")
	NotFound       = NotFoundWithMessage(404, "Not Found")
	TooManyRequest = New(429, 429, "Too Many Request")
)

func NewInvalidParamsErr(msg string) *Error {
	return BadRequest(1000, msg)
}
