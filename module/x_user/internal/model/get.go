package model

import "github.com/daidr/doulog-core/lib/models"

type GetUserInfoResp struct {
	Id        uint64 `json:"id"`
	Name      string `json:"name"`
	EmailHash string `json:"email_hash"`
	Email     string `json:"email"`
	Homepage  string `json:"homepage"`
	IsAdmin   bool   `json:"is_admin"`
	IsBanned  bool   `json:"is_banned"`
	CreatedAt int64  `json:"created_at"`
}

type UserListReq struct {
	Keyword string `form:"keyword" json:"keyword"`
	models.PageDto
}

type UserListResp struct {
	Total int64             `json:"total"`
	List  []GetUserInfoResp `json:"list"`
}
