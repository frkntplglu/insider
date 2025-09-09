package handler

type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

type FailureResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}
