package request

type EmailLoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
