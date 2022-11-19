package main

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/Luzifer/mercedes-byocar-exporter/internal/credential"
	"github.com/Luzifer/mercedes-byocar-exporter/internal/mercedes"
)

func getAuthRedirectHandler(mc mercedes.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, mc.GetAuthStartURL(cfg.RedirectURL), http.StatusTemporaryRedirect)
	}
}

func getAuthStoreTokenHandler(mc mercedes.Client, creds credential.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := mc.StoreTokenFromRequest(cfg.RedirectURL, r); err != nil {
			http.Error(w, errors.Wrap(err, "storing auth token").Error(), http.StatusInternalServerError)
			return
		}

		http.Error(w, "Token stored, configuration done.", http.StatusOK)
	}
}
