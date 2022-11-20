package influxdb

import (
	"net/url"
	"strings"
	"sync"
	"time"

	influx "github.com/influxdata/influxdb1-client/v2"
	"github.com/pkg/errors"

	"github.com/Luzifer/mercedes-byocar-exporter/internal/exporters"
)

const (
	influxTimeout       = 2 * time.Second
	influxWriteInterval = 10 * time.Second
)

type (
	Exporter struct {
		batch     influx.BatchPoints
		batchLock sync.Mutex
		client    influx.Client
		database  string
		errs      chan error
	}
)

var _ exporters.Exporter = (*Exporter)(nil)

func New(connURL string) (*Exporter, error) {
	out := &Exporter{
		errs: make(chan error, 10), //nolint: gomnd // Is a constant but makes no sense to name
	}
	return out, out.initialize(connURL)
}

func (e *Exporter) Errors() <-chan error {
	return e.errs
}

func (e *Exporter) RecordPoint(name string, tags map[string]string, fields map[string]interface{}, updatedAt time.Time) error {
	pt, err := influx.NewPoint(name, tags, fields, updatedAt)
	if err != nil {
		return err
	}

	e.batchLock.Lock()
	defer e.batchLock.Unlock()
	e.batch.AddPoint(pt)

	return nil
}

func (e *Exporter) resetBatch() error {
	b, err := influx.NewBatchPoints(influx.BatchPointsConfig{
		Database: e.database,
	})
	if err != nil {
		return err
	}

	e.batch = b
	return nil
}

func (e *Exporter) sendLoop() {
	for range time.Tick(influxWriteInterval) {

		e.batchLock.Lock()
		if err := e.client.Write(e.batch); err != nil {
			e.errs <- err
			e.batchLock.Unlock()
			continue
		}
		e.resetBatch()
		e.batchLock.Unlock()

	}
}

func (e *Exporter) initialize(connURL string) error {
	connInfo, err := url.Parse(connURL)
	if err != nil {
		return errors.Wrap(err, "parsing connection URL")
	}
	e.database = strings.TrimLeft(connInfo.Path, "/")

	cfg := influx.HTTPConfig{
		Addr:    (&url.URL{Scheme: connInfo.Scheme, Host: connInfo.Host}).String(),
		Timeout: influxTimeout,
	}

	if connInfo.User != nil {
		cfg.Username = connInfo.User.Username()
		cfg.Password = func(pass string, _ bool) string { return pass }(connInfo.User.Password())
	}

	influxClient, err := influx.NewHTTPClient(cfg)
	if err != nil {
		return err
	}

	e.client = influxClient
	if err := e.resetBatch(); err != nil {
		return err
	}
	go e.sendLoop()

	return nil
}
