package requests

import (
	"html"
	"strings"
)

type SignupRequest struct {
	Name     string `form:"name" json:"name" xml:"name" validate:"required"`
	Email    string `form:"email" json:"email" xml:"email" validate:"required,email"`
	Password string `form:"password" json:"password" xml:"password" validate:"required"`
}

func (requestData *SignupRequest) ValidateData() error {
	requestData.Name = html.EscapeString(strings.Trim(requestData.Name, " "))
	requestData.Email = html.EscapeString(strings.Trim(requestData.Email, " "))
	requestData.Password = html.EscapeString(strings.Trim(requestData.Password, " "))
	return nil
}
