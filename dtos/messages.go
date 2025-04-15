package dtos

type (
	EmailMessage struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Message  string `json:"message"`
	}
)
