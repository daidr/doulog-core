package dto

import "github.com/go-webauthn/webauthn/protocol"

type WebAuthnRegFinishReq struct {
	Challenge    string                              `json:"challenge"`
	CreationData protocol.CredentialCreationResponse `json:"data"`
}

type WebAuthnCredentialResp struct {
	ID         uint64 `json:"id"`
	Label      string `json:"label"`
	CreatedAt  int64  `json:"createdAt"`
	LastUsedAt int64  `json:"lastUsedAt"`
	Synced     bool   `json:"synced"`
}
