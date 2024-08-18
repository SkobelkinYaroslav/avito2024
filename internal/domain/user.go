package domain

type User struct {
	ID       string `json:"id,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password"`
	UserType string `json:"user_type,omitempty"`
	Token    string `json:"token"`
}
