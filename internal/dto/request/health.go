package request

type HealthReq struct {
	Message string `query:"message" validate:"required"`
}
