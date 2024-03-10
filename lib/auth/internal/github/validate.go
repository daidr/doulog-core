package github

import (
	"github.com/daidr/doulog-core/lib/conf"
	"github.com/daidr/doulog-core/lib/models"
	"github.com/daidr/doulog-core/lib/request"
	"github.com/guonaihong/gout"
	"github.com/pkg/errors"
	"strconv"
)

type AuthAccept struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	ErrorUri         string `json:"error_uri"`
	AccessToken      string `json:"access_token"`
}

type UserAccept struct {
	Login   string `json:"login"`
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Blog    string `json:"blog"`
	HtmlUrl string `json:"html_url"`
	Message string `json:"message"`
}

func Validate(code string, state string) (*models.OauthPayload, error) {
	var (
		resp     AuthAccept
		userResp UserAccept
		token    string
	)
	err := request.HTTP().POST("https://github.com/login/oauth/access_token").
		SetJSON(gout.H{
			"client_id":     conf.C.Auth.GitHub.ClientID,
			"client_secret": conf.C.Auth.GitHub.ClientSecret,
			"code":          code,
			"state":         state,
		}).
		SetHeader(gout.H{
			"accept": "application/json",
		}).
		BindJSON(&resp).Do()

	if err != nil {
		return nil, err
	}
	if resp.Error != "" {
		return nil, errors.New("wrong code")
	}

	token = resp.AccessToken

	err = request.HTTP().GET("https://api.github.com/user").
		SetHeader(gout.H{
			"Authorization": "Bearer " + token,
			"Accept":        "application/vnd.github+json",
		}).BindJSON(&userResp).Do()
	if err != nil {
		return nil, err
	}
	if userResp.Message != "" {
		return nil, errors.New("wrong code")
	}

	id := strconv.Itoa(userResp.Id)
	login := userResp.Login
	username := userResp.Name
	email := userResp.Email
	homepage := userResp.Blog
	if homepage == "" {
		homepage = userResp.HtmlUrl
	}

	return &models.OauthPayload{
		Id:       id,
		Login:    login,
		Name:     username,
		Email:    email,
		Homepage: homepage,
	}, nil
}
