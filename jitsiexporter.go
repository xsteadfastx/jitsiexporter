package jitsiexporter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

//go:generate mockery -name Stater -inpkg

type Metric struct {
	Name  string
	Gauge prometheus.Gauge
}

type Metrics struct {
	Metrics map[string]Metric
	URL     string
	Stater  Stater
	mux     sync.Mutex
}

func (m *Metrics) Update() {
	now := m.Stater.Now(m.URL)
	log.Debug(now)

	m.mux.Lock()
	for k, v := range now {
		fieldLogger := log.WithFields(log.Fields{"key": k})

		// skipping anything else than float64.
		switch t := v.(type) {
		case float64:
			fieldLogger.Debugf("found '%v'", t)

			name := fmt.Sprintf("jitsi_%s", k)
			if _, ok := m.Metrics[name]; !ok {
				fieldLogger.Info("creating and registering metric")

				m.Metrics[name] = Metric{
					Name: name,
					Gauge: prometheus.NewGauge(
						prometheus.GaugeOpts{
							Name: name,
						},
					),
				}
				fieldLogger.Debugf("%+v", m.Metrics[name])
				prometheus.MustRegister(m.Metrics[name].Gauge)
			}

			value := v.(float64)
			fieldLogger.Infof("set to %f", value)
			m.Metrics[name].Gauge.Set(value)
		default:
			fieldLogger.Debugf("found %v", t)
			fieldLogger.Info("skipping")

			continue
		}
	}
	m.mux.Unlock()
}

type Stater interface {
	Now(url string) map[string]interface{}
}

type colibri struct{}

func (c colibri) Now(url string) map[string]interface{} {
	s := make(map[string]interface{})
	resp, err := http.Get(url) // nolint:gosec

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&s)

	if err != nil {
		log.Fatal(err)
	}

	return s
}

func collect(m *Metrics) {
	for {
		m.Update()
		time.Sleep(30 * time.Second) // nolint:gomnd
	}
}

func Serve(url string, debug bool, interval time.Duration, port int, host string) {
	s := colibri{}
	metrics := &Metrics{
		URL:     url,
		Stater:  s,
		Metrics: make(map[string]Metric),
	}

	if debug {
		log.SetLevel(log.DebugLevel)
	}

	log.Debugf("%+v", metrics)

	go collect(metrics)

	http.Handle("/metrics", promhttp.Handler())
	log.Info("beginning to serve")
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), nil))
}
