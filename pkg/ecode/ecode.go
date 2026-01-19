package ecode

var (
	InternalServerErr = InternalServer(-1, "Internal Server Error")

	InvalidParams  = NewInvalidParamsErr("Invalid Params")
	NotFound       = NewInvalidParamsErr("Not Found")
	TooManyRequest = BadRequest(1100, "Too Many Request")

	InvalidToken = BadRequest(2000, "Invalid Token")
	ExpiresToken = BadRequest(2001, "Token Expires")

	EmailOrPasswordErr = BadRequest(2100, "email or password error")
)

func NewInvalidParamsErr(msg string) *Error {
	return BadRequest(1000, msg)
}
