package credential

import (
	"encoding/json"
	"io/fs"
	"os"
	"time"

	"github.com/pkg/errors"
)

type (
	JSONStore struct {
		clientID, clientSecret, filename string
	}

	jsonContent struct {
		AccessToken, RefreshToken string
		Expiry                    time.Time
	}
)

var _ Store = JSONStore{}

func NewJSONStore(filename, clientID, clientSecret string) (JSONStore, error) {
	store := JSONStore{
		clientID:     clientID,
		clientSecret: clientSecret,
		filename:     filename,
	}

	_, err := os.Stat(filename)
	switch {
	case errors.Is(err, nil), errors.Is(err, fs.ErrNotExist):
		return store, nil

	default:
		return store, errors.Wrap(err, "probing store file")
	}
}

func (j JSONStore) GetClientCredentials() (clientID, clientSecret string, err error) {
	return j.clientID, j.clientSecret, nil
}

func (j JSONStore) GetToken() (accessToken, refreshToken string, expiry time.Time, err error) {
	var c jsonContent

	f, err := os.Open(j.filename)
	if err != nil {
		return "", "", time.Time{}, errors.Wrap(err, "opening store")
	}
	defer f.Close()

	if err = json.NewDecoder(f).Decode(&c); err != nil {
		return "", "", time.Time{}, errors.Wrap(err, "decoding store")
	}

	return c.AccessToken, c.RefreshToken, c.Expiry, nil
}

func (j JSONStore) HasCredentials() (bool, error) {
	_, r, _, err := j.GetToken()
	switch {
	case errors.Is(err, fs.ErrNotExist):
		return false, nil

	case errors.Is(err, nil):
		return r != "", nil

	default:
		return false, errors.Wrap(err, "getting credentials")
	}
}

func (j JSONStore) UpdateToken(accessToken, refreshToken string, expiry time.Time) error {
	c := jsonContent{accessToken, refreshToken, expiry}

	f, err := os.Create(j.filename)
	if err != nil {
		return errors.Wrap(err, "creating store file")
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(c)
}
