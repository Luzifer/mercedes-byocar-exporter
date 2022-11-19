package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"

	"github.com/Luzifer/mercedes-byocar-exporter/internal/credential"
	"github.com/Luzifer/mercedes-byocar-exporter/internal/mercedes"
	"github.com/Luzifer/rconfig/v2"
)

var (
	cfg     cliConfig
	version = "dev"
)

func initApp() error {
	rconfig.AutoEnv(true)
	if err := rconfig.ParseAndValidate(&cfg); err != nil {
		return errors.Wrap(err, "parsing cli options")
	}

	if cfg.VersionAndExit {
		fmt.Printf("mercedes-byocar-exporter %s\n", version)
		os.Exit(0)
	}

	l, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		return errors.Wrap(err, "parsing log-level")
	}
	logrus.SetLevel(l)

	if err = cfg.Validate(); err != nil {
		return errors.Wrap(err, "validating config")
	}

	return nil
}

func main() {
	var err error
	if err = initApp(); err != nil {
		logrus.WithError(err).Fatal("initializing app")
	}

	// Initialize credentials store
	var creds credential.Store
	switch {
	case cfg.ClientID != "":
		logrus.WithField("method", "json-file").Debug("opening credential store")
		creds, err = credential.NewJSONStore(cfg.CredentialFile, cfg.ClientID, cfg.ClientSecret)
	case cfg.VaultKey != "":
		logrus.WithField("method", "vault").Debug("opening credential store")
		creds, err = credential.NewVaultStore(cfg.VaultKey)
	}
	if err != nil {
		logrus.WithError(err).Fatal("initializing credential store")
	}
	logrus.WithField("method", "vault").Debug("credential store connected")

	// Initialize Mercedes API client
	clientID, clientSecret, err := creds.GetClientCredentials()
	if err != nil {
		logrus.WithError(err).Fatal("getting client credentials")
	}
	mClient := mercedes.New(clientID, clientSecret, creds)

	// Register HTTP handlers
	http.DefaultServeMux.HandleFunc("/auth", getAuthRedirectHandler(mClient))
	http.DefaultServeMux.HandleFunc("/store-token", getAuthStoreTokenHandler(mClient, creds))
	http.DefaultServeMux.Handle("/metrics", promhttp.Handler())

	scheduler := cron.New()
	scheduler.AddFunc(fmt.Sprintf("@every %s", cfg.FetchInterval), getCronFunc(mClient))
	scheduler.Start()

	// Do an initial fetch to propagate metrics
	getCronFunc(mClient)()

	// Start HTTP server
	logrus.WithField("version", version).Info("mercedes-byocar-exporter started")
	srv := http.Server{
		Addr:              cfg.Listen,
		Handler:           http.DefaultServeMux,
		ReadHeaderTimeout: time.Second,
	}
	if err = srv.ListenAndServe(); err != nil {
		logrus.WithError(err).Fatal("HTTP server exitted unexpectedly")
	}
}
