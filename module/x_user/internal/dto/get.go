package dto

import "github.com/daidr/doulog-core/lib/models"

type GetUserInfoResp struct {
	Id        uint64 `json:"id"`
	Name      string `json:"name"`
	EmailHash string `json:"emailHash"`
	Email     string `json:"email"`
	Homepage  string `json:"homepage"`
	Motto     string `json:"motto"`
	IsAdmin   bool   `json:"isAdmin"`
	IsBanned  bool   `json:"isBanned"`
	CreatedAt int64  `json:"createdAt"`
}

type UserListReq struct {
	Keyword string `form:"keyword" json:"keyword"`
	models.PageDto
}

type UserListResp struct {
	Total int64             `json:"total"`
	List  []GetUserInfoResp `json:"list"`
}
