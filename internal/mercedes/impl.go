package mercedes

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"

	"github.com/Luzifer/mercedes-byocar-exporter/internal/credential"
)

const (
	requestTimeout  = 10 * time.Second
	stateExpiry     = 5 * time.Minute
	tokenGraceRenew = -5 * time.Minute
)

type (
	APIClient struct {
		clientID, clientSecret string
		creds                  credential.Store

		validStateToken       string
		validStateTokenExpiry time.Time
	}
)

var _ Client = (*APIClient)(nil)

func New(clientID, clientSecret string, creds credential.Store) *APIClient {
	return &APIClient{
		clientID:     clientID,
		clientSecret: clientSecret,
		creds:        creds,
	}
}

func (a *APIClient) GetAuthStartURL(redirectURL string) string {
	a.validStateToken = uuid.Must(uuid.NewV4()).String()
	a.validStateTokenExpiry = time.Now().Add(stateExpiry)

	return a.getOauth2Config(redirectURL).AuthCodeURL(a.validStateToken, oauth2.AccessTypeOffline)
}

func (a *APIClient) StoreTokenFromRequest(redirectURL string, r *http.Request) error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	state := r.FormValue("state")
	if state != a.validStateToken || a.validStateTokenExpiry.Before(time.Now()) {
		return errors.New("invalid or expired state")
	}

	code := r.FormValue("code")
	tok, err := a.getOauth2Config(redirectURL).Exchange(ctx, code, oauth2.AccessTypeOffline)
	if err != nil {
		return errors.Wrap(err, "exchanging code for token")
	}

	return errors.Wrap(a.creds.UpdateToken(tok.AccessToken, tok.RefreshToken, tok.Expiry), "updating stored token")
}

func (a APIClient) getOauth2Config(redirectURL string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     a.clientID,
		ClientSecret: a.clientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:   oAuthEndpointAuth,
			TokenURL:  oAuthEndpointToken,
			AuthStyle: oauth2.AuthStyleInHeader,
		},
		RedirectURL: redirectURL,
		Scopes: []string{
			oAuthScopeOfflineAccess,
			oAuthScopeOpenID,
			oAuthScopePayAsYouDrive,
			oAuthScopeVehicleElectricStatus,
			oAuthScopeVehicleFuelStatus,
			oAuthScopeVehicleLockStatus,
			oAuthScopeVehicleStatus,
		},
	}
}

func (a APIClient) parseGenericAPIResponse(data io.Reader, output any) (err error) {
	var tmp genericAPIResponse
	if err = json.NewDecoder(data).Decode(&tmp); err != nil {
		return errors.Wrap(err, "parsing JSON response")
	}

	st := reflect.ValueOf(output).Elem()
	for i := 0; i < st.NumField(); i++ {
		valField := st.Field(i)
		typeField := st.Type().Field(i)

		name := typeField.Tag.Get("apiField")
		if name == "" {
			continue
		}

		value := tmp.Get(name)
		if value == nil {
			continue
		}
		value.Timestamp = value.Timestamp * 1000000

		switch typeField.Type {
		case reflect.TypeOf(TimedBool{}):
			fv := TimedBool{}
			if fv.v, err = strconv.ParseBool(value.Value); err != nil {
				return errors.Wrapf(err, "parsing value for %s", name)
			}
			fv.t = time.Unix(0, value.Timestamp)

			valField.Set(reflect.ValueOf(fv))

		case reflect.TypeOf(TimedEnum{}):
			fv := TimedEnum{}
			if fv.v, err = strconv.ParseInt(value.Value, 10, 64); err != nil {
				return errors.Wrapf(err, "parsing value for %s", name)
			}
			fv.t = time.Unix(0, value.Timestamp)
			fv.def = strings.Split(typeField.Tag.Get("values"), ",")

			valField.Set(reflect.ValueOf(fv))

		case reflect.TypeOf(TimedFloat{}):
			fv := TimedFloat{}
			if fv.v, err = strconv.ParseFloat(value.Value, 64); err != nil {
				return errors.Wrapf(err, "parsing value for %s", name)
			}
			fv.t = time.Unix(0, value.Timestamp)

			valField.Set(reflect.ValueOf(fv))

		case reflect.TypeOf(TimedInt{}):
			fv := TimedInt{}
			if fv.v, err = strconv.ParseInt(value.Value, 10, 64); err != nil {
				return errors.Wrapf(err, "parsing value for %s", name)
			}
			fv.t = time.Unix(0, value.Timestamp)

			valField.Set(reflect.ValueOf(fv))

		default:
			return errors.Errorf("unknown type %v", typeField.Type)
		}
	}

	return nil
}

func (a APIClient) request(path string, output any) error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	url := strings.Join([]string{
		strings.TrimRight(apiPrefix, "/"),
		strings.TrimLeft(path, "/"),
	}, "/")

	logrus.WithField("url", url).Trace("creating request")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return errors.Wrap(err, "creating request")
	}
	req.Header.Set("accept", "application/json;charset=utf-8")

	at, rt, exp, err := a.creds.GetToken()
	if err != nil {
		return errors.Wrap(err, "getting credentials")
	}
	tok := &oauth2.Token{AccessToken: at, RefreshToken: rt, Expiry: exp}

	// Renew token if required
	if tok.Expiry.Add(tokenGraceRenew).Before(time.Now()) {
		src := a.getOauth2Config("").TokenSource(ctx, tok)
		if tok, err = src.Token(); err != nil {
			return errors.Wrap(err, "renewing token")
		}

		if tok.AccessToken != at || tok.RefreshToken != rt || !tok.Expiry.Equal(exp) {
			if err := a.creds.UpdateToken(tok.AccessToken, tok.RefreshToken, tok.Expiry); err != nil {
				return errors.Wrap(err, "updating stored token")
			}
		}
	}

	resp, err := a.getOauth2Config("").Client(ctx, tok).Do(req)
	if err != nil {
		return errors.Wrap(err, "executing request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrapf(err, "http status code %d, error reading body", resp.StatusCode)
		}
		return errors.Errorf("http status code %d, body %s", resp.StatusCode, body)
	}

	if output == nil {
		return nil
	}

	return errors.Wrap(a.parseGenericAPIResponse(resp.Body, output), "decoding output")
}
