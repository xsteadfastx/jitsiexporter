package jitsiexporter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
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
	Exclude []string
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

		if sort.SearchStrings(m.Exclude, k) != len(m.Exclude) {
			fieldLogger.Info("exclude")

			continue
		}

		name := fmt.Sprintf("jitsi_%s", k)
		if _, ok := m.Metrics[name]; !ok {
			fieldLogger.Info("creating and registerting metric")

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

func Serve(url string) {
	s := colibri{}
	metrics := &Metrics{
		Exclude: []string{
			"conference_sizes",
			"current_timestamp",
			"graceful_shutdown",
		},
		URL:     url,
		Stater:  s,
		Metrics: make(map[string]Metric),
	}

	log.SetLevel(log.DebugLevel)
	log.Debugf("%+v", metrics)

	go collect(metrics)

	http.Handle("/metrics", promhttp.Handler())
	log.Info("beginning to serve")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
