package webauthn

import (
	"github.com/go-webauthn/webauthn/webauthn"
	"time"
)

var (
	WebAuthn *webauthn.WebAuthn
)

func Init() error {
	config := &webauthn.Config{
		RPDisplayName: "DouLog",
		RPID:          "localhost",
		RPOrigins:     []string{"http://localhost:3003"},
		Timeouts: webauthn.TimeoutsConfig{
			Login: webauthn.TimeoutConfig{
				Enforce:    true,             // Require the response from the client comes before the end of the timeout.
				Timeout:    time.Second * 60, // Standard timeout for login sessions.
				TimeoutUVD: time.Second * 60, // Timeout for login sessions which have user verification set to discourage.
			},
			Registration: webauthn.TimeoutConfig{
				Enforce:    true,             // Require the response from the client comes before the end of the timeout.
				Timeout:    time.Second * 60, // Standard timeout for registration sessions.
				TimeoutUVD: time.Second * 60, // Timeout for login sessions which have user verification set to discourage.
			},
		},
	}

	var err error
	WebAuthn, err = webauthn.New(config)
	if err != nil {
		return err
	}
	return nil
}
