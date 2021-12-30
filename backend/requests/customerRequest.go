package requests

import (
	"html"
	"strings"
)

type CustomerRequest struct {
	Name    string `form:"name" json:"name" xml:"name" validate:"required"`
	Email   *string
	Address *string
	Phone   *string
}

func (requestData *CustomerRequest) ValidateData() error {

	requestData.Name = html.EscapeString(strings.Trim(requestData.Name, " "))

	if requestData.Email != nil {
		email := html.EscapeString(strings.Trim(*requestData.Email, " "))
		requestData.Email = &email
	}

	if requestData.Address != nil {
		address := html.EscapeString(strings.Trim(*requestData.Address, " "))
		requestData.Address = &address
	}

	if requestData.Phone != nil {
		phone := html.EscapeString(strings.Trim(*requestData.Phone, " "))
		requestData.Phone = &phone
	}

	return nil
}
