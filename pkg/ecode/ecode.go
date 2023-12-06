package ecode

var (
	InvalidParams     = NewInvalidParamsErr("Invalid Params")
	NotFound          = NewInvalidParamsErr("Not Found")
	TooManyRequest    = New(1100, 429, "Too Many Request")
	InvalidHashID     = BadRequest(1101, "invalid hash id")
	InternalServerErr = InternalServer(500, "Internal Server Error")

	// user
	InvalidToken      = BadRequest(2000, "Invalid Token")
	ExpiresToken      = BadRequest(2001, "Token Expires")
	InvalidOAuthState = BadRequest(2002, "Invalid OAuth State")
)

func NewInvalidParamsErr(msg string) *Error {
	return BadRequest(1000, msg)
}
