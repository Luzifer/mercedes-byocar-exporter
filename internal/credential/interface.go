package credential

import "time"

type (
	Store interface {
		GetClientCredentials() (clientID, clientSecret string, err error)
		GetToken() (accessToken, refreshToken string, expiry time.Time, err error)
		HasCredentials() (bool, error)
		UpdateToken(accessToken, refreshToken string, expiry time.Time) error
	}
)
