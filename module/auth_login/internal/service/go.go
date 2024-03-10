package service

import (
	"context"
	"time"

	"github.com/daidr/doulog-core/lib/auth"
	"github.com/daidr/doulog-core/lib/format"
	"github.com/daidr/doulog-core/lib/models"
	"github.com/pkg/errors"
)

func Go(db *models.DB, platform string, frontendCallback string) (string, error) {
	redirect, mark := auth.GetRedirectUrl(platform)
	if redirect == "" {
		return "", errors.New("can not get redirect url")
	}

	// mark 15分钟
	err := db.Redis.Set(context.Background(), format.Key.AuthCallback(mark), frontendCallback, time.Minute*15).Err()
	if err != nil {
		return "", err
	}
	return redirect, nil
}
