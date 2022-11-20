package credential

import (
	"os"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type (
	VaultStore struct {
		client *api.Client
		key    string
	}
)

var _ Store = VaultStore{}

var ErrMissingKey = errors.New("missing key")

func NewVaultStore(vaultKey string) (VaultStore, error) {
	var (
		err error
		v   VaultStore
	)

	v.key = vaultKey

	c := api.DefaultConfig()
	if err = c.ReadEnvironment(); err != nil {
		return v, errors.Wrap(err, "configuring Vault from ENV")
	}

	v.client, err = api.NewClient(c)
	if err != nil {
		return v, errors.Wrap(err, "creating Vault client")
	}

	return v, errors.Wrap(v.authorizeVault(), "authorizing Vault")
}

func (v VaultStore) GetClientCredentials() (clientID, clientSecret string, err error) {
	if err = v.authorizeVault(); err != nil {
		return "", "", errors.Wrap(err, "authorizing Vault")
	}

	secret, err := v.client.Logical().Read(v.key)
	if err != nil {
		return "", "", errors.Wrap(err, "reading Vault key")
	}
	if secret == nil || secret.Data == nil {
		logrus.Warnf("%s %v", v.key, secret)
		return "", "", errors.New("no data found at key")
	}

	var ok bool
	if clientID, ok = secret.Data["client-id"].(string); !ok {
		return "", "", errors.Wrap(ErrMissingKey, "getting client-id")
	}
	if clientSecret, ok = secret.Data["client-secret"].(string); !ok {
		return "", "", errors.Wrap(ErrMissingKey, "getting client-secret")
	}

	return clientID, clientSecret, nil
}

func (v VaultStore) GetToken() (accessToken, refreshToken string, expiry time.Time, err error) {
	if err = v.authorizeVault(); err != nil {
		return "", "", time.Time{}, errors.Wrap(err, "authorizing Vault")
	}

	secret, err := v.client.Logical().Read(v.key)
	if err != nil {
		return "", "", time.Time{}, errors.Wrap(err, "reading Vault key")
	}
	if secret == nil || secret.Data == nil {
		return "", "", time.Time{}, errors.New("no data found at key")
	}

	var ok bool
	if accessToken, ok = secret.Data["access-token"].(string); !ok {
		return "", "", time.Time{}, errors.Wrap(ErrMissingKey, "getting access-token")
	}
	if refreshToken, ok = secret.Data["refresh-token"].(string); !ok {
		return "", "", time.Time{}, errors.Wrap(ErrMissingKey, "getting refresh-token")
	}
	exp, ok := secret.Data["expiry"].(string)
	if !ok {
		return "", "", time.Time{}, errors.Wrap(ErrMissingKey, "getting expiry")
	}
	if expiry, err = time.Parse(time.RFC3339Nano, exp); err != nil {
		return "", "", time.Time{}, errors.Wrap(err, "parsing stored expiry")
	}

	return accessToken, refreshToken, expiry, nil
}

func (v VaultStore) HasCredentials() (bool, error) {
	_, r, _, err := v.GetToken()
	switch {
	case errors.Is(err, nil):
		return r != "", nil

	case errors.Is(err, ErrMissingKey):
		return false, nil

	default:
		return false, errors.Wrap(err, "getting client credentials")
	}
}

func (v VaultStore) UpdateToken(accessToken, refreshToken string, expiry time.Time) (err error) {
	if err = v.authorizeVault(); err != nil {
		return errors.Wrap(err, "authorizing Vault")
	}

	secret, err := v.client.Logical().Read(v.key)
	if err != nil {
		return errors.Wrap(err, "reading Vault key")
	}
	if secret == nil || secret.Data == nil {
		return errors.New("no data found at key")
	}

	data := secret.Data
	data["access-token"] = accessToken
	data["refresh-token"] = refreshToken
	data["expiry"] = expiry.Format(time.RFC3339Nano)

	_, err = v.client.Logical().Write(v.key, data)
	return errors.Wrap(err, "writing back data")
}

func (v VaultStore) authorizeVault() error {
	if role := os.Getenv("VAULT_ROLE_ID"); role != "" {
		data := map[string]interface{}{
			"role_id": role,
		}
		if secret := os.Getenv("VAULT_SECRET_ID"); secret != "" {
			data["secret_id"] = secret
		}
		loginSecret, lserr := v.client.Logical().Write("auth/approle/login", data)
		if lserr != nil || loginSecret.Auth == nil {
			return errors.Wrap(lserr, "fetching token for approle")
		}

		v.client.SetToken(loginSecret.Auth.ClientToken)
		return nil
	}

	if token := os.Getenv(api.EnvVaultToken); token != "" {
		v.client.SetToken(token)
		return nil
	}

	if tokenFile, err := homedir.Expand("~/.vault-token"); err == nil {
		if token, err := os.ReadFile(tokenFile); err == nil {
			v.client.SetToken(string(token))
			return nil
		}
	}

	return errors.New("no valid auth method found")
}
