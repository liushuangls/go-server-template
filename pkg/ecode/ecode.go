package ecode

var (
	InvalidParams  = NewInvalidParamsErr("Invalid Params")
	NotFound       = NewInvalidParamsErr("Not Found")
	TooManyRequest = BadRequest(1100, "Too Many Request")

	// user
	InvalidToken      = BadRequest(2000, "Invalid Token")
	ExpiresToken      = BadRequest(2001, "Token Expires")
	InvalidOAuthState = BadRequest(2002, "Invalid OAuth State")
)

func NewInvalidParamsErr(msg string) *Error {
	return BadRequest(1000, msg)
}
