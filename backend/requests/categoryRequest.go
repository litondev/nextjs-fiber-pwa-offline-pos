package requests

import (
	"html"
	"strings"
)

type CategoryRequest struct {
	Name        string `form:"name" json:"name" xml:"name" validate:"required"`
	Description *string
}

func (requestData *CategoryRequest) ValidateData() error {
	requestData.Name = html.EscapeString(strings.Trim(requestData.Name, " "))
	if requestData.Description != nil {
		des := html.EscapeString(strings.Trim(*requestData.Description, " "))
		requestData.Description = &des
	}
	return nil
}
