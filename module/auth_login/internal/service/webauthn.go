package service

import (
	"context"
	"encoding/json"
	"github.com/daidr/doulog-core/lib/daos"
	"github.com/daidr/doulog-core/lib/format"
	"github.com/daidr/doulog-core/lib/models"
	"github.com/daidr/doulog-core/lib/webauthn"
	"github.com/daidr/doulog-core/module/auth_login/internal/dto"
	authUtils "github.com/daidr/doulog-core/module/auth_login/internal/utils"
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
	u, err := daos.NewUser(sp.DB).GetCredentials(uid)
	if err != nil {
		return nil, err
	}
	exclusions := make([]protocol.CredentialDescriptor, 0, len(u.Credentials))
	for _, c := range u.Credentials {
		exclusions = append(exclusions, protocol.CredentialDescriptor{
			Type:         protocol.PublicKeyCredentialType,
			CredentialID: c.Credential.Data().ID,
		})
	}
	options, session, err := webauthn.WebAuthn.BeginRegistration(u, webauthn2.WithExclusions(exclusions), webauthn2.WithResidentKeyRequirement(protocol.ResidentKeyRequirementRequired))
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
	u, err := daos.NewUser(sp.DB).GetCredentials(uid)
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
	u, err := daos.NewUser(sp.DB).GetCredentials(uid)
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

func DeleteWebAuthnCredential(sp *models.Scope, uid, cid uint64) error {
	err := daos.NewUser(sp.DB).DeleteCredential(uid, cid)

	return err
}

func RenameWebAuthnCredential(sp *models.Scope, uid, cid uint64, newName string) error {
	err := daos.NewUser(sp.DB).RenameCredential(uid, cid, newName)

	return err
}

func BeginWebAuthnDiscoverLoginOptions(sp *models.Scope) (*protocol.CredentialAssertion, error) {
	options, session, err := webauthn.WebAuthn.BeginDiscoverableLogin()
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

func FinishWebAuthnDiscoverLogin(sp *models.Scope, req dto.WebAuthnLoginFinishReq) (string, error) {
	sessionJSON, err := sp.DB.Redis.Get(context.Background(), format.Key.WebauthnSession(req.Challenge)).Result()
	if err != nil {
		return "", err
	}

	session := &customSessionData{}
	if err = json.Unmarshal([]byte(sessionJSON), session); err != nil {
		return "", err
	}

	parsedData, err := req.AssertionData.Parse()
	if err != nil {
		return "", err
	}

	credential, err := webauthn.WebAuthn.ValidateDiscoverableLogin(func(rawID, userHandle []byte) (user webauthn2.User, err error) {
		return daos.NewUser(sp.DB).GetUserByCredential(rawID, userHandle)
	}, session.SessionData, parsedData)
	if err != nil {
		return "", err
	}

	uid := models.WebAuthnIDToUint64(parsedData.Response.UserHandle)

	err = daos.NewUser(sp.DB).UpdateCredentialLastUsedAt(models.WebAuthnIDToUint64(parsedData.Response.UserHandle), credential)

	token := authUtils.SetToken(sp.DB, uid)
	return token, err
}
