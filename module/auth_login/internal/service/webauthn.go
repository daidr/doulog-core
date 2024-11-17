package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/daidr/doulog-core/lib/daos"
	"github.com/daidr/doulog-core/lib/format"
	"github.com/daidr/doulog-core/lib/models"
	"github.com/daidr/doulog-core/lib/webauthn"
	"github.com/daidr/doulog-core/module/auth_login/internal/dto"
	"github.com/go-webauthn/webauthn/protocol"
	webauthn2 "github.com/go-webauthn/webauthn/webauthn"
	"time"
)

type customSessionData struct {
	webauthn2.SessionData
}

func (i customSessionData) MarshalBinary() ([]byte, error) {
	return json.Marshal(i)
}

func (i customSessionData) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &i)
}

func BeginWebAuthnRegOptions(sp *models.Scope, uid uint64) (*protocol.CredentialCreation, error) {
	u, err := daos.NewUser(sp.DB).GetWithCredentials(uid)
	if err != nil {
		return nil, err
	}
	// format printout of u.WebAuthnCredentials() (a array)
	fmt.Println(u.WebAuthnCredentials())
	options, session, err := webauthn.WebAuthn.BeginRegistration(u)
	if err != nil {
		return nil, err
	}

	codableSession := &customSessionData{SessionData: *session}

	sessionJSON, err := json.Marshal(codableSession)
	if err != nil {
		return nil, err
	}

	if err = sp.DB.Redis.Set(context.Background(), format.Key.WebauthnSession(session.Challenge), sessionJSON, time.Until(session.Expires)).Err(); err != nil {
		return nil, err
	}

	return options, nil
}

func FinishWebAuthnReg(sp *models.Scope, uid uint64, req dto.WebAuthnRegFinishReq) (credential *webauthn2.Credential, err error) {
	u, err := daos.NewUser(sp.DB).GetWithCredentials(uid)
	if err != nil {
		return nil, err
	}

	sessionJSON, err := sp.DB.Redis.Get(context.Background(), format.Key.WebauthnSession(req.Challenge)).Result()
	if err != nil {
		return nil, err
	}

	session := &customSessionData{}
	if err = json.Unmarshal([]byte(sessionJSON), session); err != nil {
		return nil, err
	}

	parsedData, err := req.CreationData.Parse()
	if err != nil {
		return nil, err
	}

	credential, err = webauthn.WebAuthn.CreateCredential(u, session.SessionData, parsedData)
	return credential, err
}

func ListWebAuthnCredentials(sp *models.Scope, uid uint64) ([]dto.WebAuthnCredentialResp, error) {
	u, err := daos.NewUser(sp.DB).GetWithCredentials(uid)
	if err != nil {
		return nil, err
	}
	rawCredentials := u.Credentials
	credentials := make([]dto.WebAuthnCredentialResp, 0, len(rawCredentials))
	for _, c := range rawCredentials {
		credentials = append(credentials, dto.WebAuthnCredentialResp{
			ID:         c.ID,
			Label:      c.Label,
			LastUsedAt: c.LastUsedAt,
			CreatedAt:  c.CreatedAt,
			Synced:     c.Credential.Data().Flags.BackupEligible && c.Credential.Data().Flags.BackupState,
		})
	}

	return credentials, nil
}
