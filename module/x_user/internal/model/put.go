package model

type UpdateUserInfoReq struct {
	NewName     string `json:"name" form:"name" binding:"omitempty,min=1,max=25,sensitive,xss"`
	NewEmail    string `json:"email" form:"email" binding:"omitempty,email"`
	NewHomepage string `json:"homepage" form:"homepage" binding:"omitempty,min=1,max=99,http_url"`
}

type UpdateUserInfoUri struct {
	TargetUid   uint64 `uri:"uid"`
	TargetField string `uri:"field" binding:"required,oneof=name email homepage"`
}
