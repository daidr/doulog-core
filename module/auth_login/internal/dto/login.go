package dto

type PwLoginReq struct {
	Email    string `json:"email" form:"email" binding:"required,email"`
	Password string `json:"password" form:"password" binding:"required"`
}

type PwLoginResp struct {
	Token string `json:"token"`
}
