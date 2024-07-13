package responses

type ErrResponse struct {
	Error string
}

type SuccessResponse struct {
	Message string `json:"message"`
}
