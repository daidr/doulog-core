package model

type GetSelfInfoResp struct {
	Id        uint64 `json:"id"`
	Name      string `json:"name"`
	EmailHash string `json:"email_hash"`
	Email     string `json:"email"`
	IsAdmin   bool   `json:"is_admin"`
	IsBanned  bool   `json:"is_banned"`
}
