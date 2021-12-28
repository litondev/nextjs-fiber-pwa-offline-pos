package requests

type SigninRequest struct {
	Email string `form:"email" json:"email" xml:"email" validate:"required,email"`
 	Password string `form:"password" json:"password" xml:"password" validate:"required"`
}

func (m SigninRequest) Structure() string {
	return "test"
}