package requests

import (
	"html"
	"strings"
)

type SigninRequest struct {
	Email    string `form:"email" json:"email" xml:"email" validate:"required,email"`
	Password string `form:"password" json:"password" xml:"password" validate:"required"`
}

func (requestData *SigninRequest) ValidateData() error {
	requestData.Email = html.EscapeString(strings.Trim(requestData.Email, " "))
	requestData.Password = html.EscapeString(strings.Trim(requestData.Password, " "))
	return nil
}
