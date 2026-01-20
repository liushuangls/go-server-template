package ecode

var (
	InternalServerErr = InternalServer(-1, "Internal Server Error")

	InvalidParams  = NewInvalidParamsErr("Invalid Params")
	NotFound       = NotFoundWithMessage(1000, "Not Found")
	TooManyRequest = New(1100, 429, "Too Many Request")
)

func NewInvalidParamsErr(msg string) *Error {
	return BadRequest(1000, msg)
}
