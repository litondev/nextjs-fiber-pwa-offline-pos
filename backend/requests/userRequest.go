package requests

import (
	"html"
	"strings"
)

type UserRequest struct {
	Name     string `form:"name" json:"name" xml:"name" validate:"required"`
	Email    string `form:"email" json:"email" xml:"email" validate:"required,email"`
	Password *string
}

func (requestData *UserRequest) ValidateData() error {

	requestData.Name = html.EscapeString(strings.Trim(requestData.Name, " "))
	requestData.Email = html.EscapeString(strings.Trim(requestData.Email, " "))

	if requestData.Password != nil {
		password := html.EscapeString(strings.Trim(*requestData.Password, " "))
		requestData.Password = &password
	}

	return nil
}
