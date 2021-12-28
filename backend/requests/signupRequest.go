package requests

type SignupRequest struct {
	Name string `form:"name" json:"name" xml:"name" validate:"required"`
	Email string `form:"email" json:"email" xml:"email" validate:"required,email"`
 	Password string `form:"password" json:"password" xml:"password" validate:"required"`
}

func (m SignupRequest) Structure() string {
	return "test"
}