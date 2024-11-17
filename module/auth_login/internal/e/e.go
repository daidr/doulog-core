// Package e 错误码定义
package e

import (
	"github.com/daidr/doulog-core/lib/ecode"
)

var (
	ErrInvalidCallback         = ecode.New(101001)
	ErrUserExists              = ecode.New(101002)
	ErrEmailNotExists          = ecode.New(101003)
	ErrPasswordWrong           = ecode.New(101004)
	ErrCreateWebAuthnChallenge = ecode.New(101005)
	ErrFinishWebAuthnVerify    = ecode.New(101006)
)

var ECode = map[ecode.Code]string{
	ErrInvalidCallback:         "Invalid callback",
	ErrUserExists:              "User already exists",
	ErrEmailNotExists:          "Email not exists",
	ErrPasswordWrong:           "Password error",
	ErrCreateWebAuthnChallenge: "Failed to create webauthn challenge",
	ErrFinishWebAuthnVerify:    "Failed to finish webauthn verify",
}
