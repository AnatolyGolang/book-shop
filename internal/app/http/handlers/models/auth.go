package models

import "fmt"

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *AuthRequest) Validate() error {
	if r.Email == "" {
		return fmt.Errorf("email is required")
	}
	if r.Password == "" {
		return fmt.Errorf("password is required")
	}
	return nil
}
