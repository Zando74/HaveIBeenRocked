package http_controller

type CheckPasswordValidator struct {
	Password string `json:"password" validate:"required"`
}

type CheckPasswordResponse struct {
	Found bool `json:"found"`
}
