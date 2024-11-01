package modelresponses

type ErrorResponse struct {
	Message string `json:"message"`
}

func ToErrorResponse(message string) ErrorResponse {
	return ErrorResponse{
		Message: message,
	}
}
