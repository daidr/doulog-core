package dto

type RegisterReq struct {
	Email    string `json:"email" form:"email" binding:"required,email"`
	Name     string `json:"name" form:"name" binding:"required,min=1,max=25"`
	Password string `json:"password" form:"password" binding:"required,min=8,max=25"`
}
