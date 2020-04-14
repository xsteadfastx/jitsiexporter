package jitsiexporter

import (
	"context"
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
	Gauge *prometheus.GaugeVec
}

type Metrics struct {
	Metrics  map[string]Metric
	URL      string
	Stater   Stater
	mux      sync.Mutex
	Errors   prometheus.Counter
	Interval time.Duration
	Hostname string
}

func (m *Metrics) Update() error {
	m.mux.Lock()
	defer m.mux.Unlock()

	now, err := m.Stater.Now(m.URL)

	if err != nil {
		m.Errors.Inc()

		for k, v := range m.Metrics {
			prometheus.Unregister(v.Gauge)
			delete(m.Metrics, k)
		}

		return err
	}

	log.Debug(now)

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
					Gauge: prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: name,
						},
						[]string{"hostname"},
					),
				}
				fieldLogger.Debugf("%+v", m.Metrics[name])
				prometheus.MustRegister(m.Metrics[name].Gauge)
			}

			value := v.(float64)
			fieldLogger.Infof("set to %f", value)
			m.Metrics[name].Gauge.WithLabelValues(m.Hostname).Set(value)
		default:
			fieldLogger.Debugf("found %v", t)
			fieldLogger.Info("skipping")

			continue
		}
	}

	return nil
}

type Response struct {
	Resp  *http.Response
	Error error
}

func get(ctx context.Context, url string, resp chan Response) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		resp <- Response{Resp: nil, Error: err}
		return
	}

	client := http.DefaultClient

	res, err := client.Do(req.WithContext(ctx)) // nolint:bodyclose
	if err != nil {
		resp <- Response{Resp: nil, Error: err}
	}

	resp <- Response{Resp: res, Error: nil}
}

type Stater interface {
	Now(url string) (map[string]interface{}, error)
}

type colibri struct{}

func (c colibri) Now(url string) (map[string]interface{}, error) {
	s := make(map[string]interface{})
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second) // nolint:gomnd

	defer cancel()

	res := make(chan Response)

	var resp *http.Response

	var err error

	go get(ctx, url, res)

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case r := <-res:
		err = r.Error
		resp = r.Resp

		defer resp.Body.Close()
	}

	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(resp.Body).Decode(&s)

	if err != nil {
		return nil, err
	}

	return s, nil
}

func collect(m *Metrics) {
	for {
		err := m.Update()
		if err != nil {
			log.Error(err)
		}

		time.Sleep(m.Interval) // nolint:gomnd
	}
}

func Serve(url string, debug bool, interval time.Duration, port int, host string, servername string) {
	s := colibri{}
	e := prometheus.NewCounter(prometheus.CounterOpts{Name: "jitsi_fetch_errors"})
	metrics := &Metrics{
		URL:      url,
		Stater:   s,
		Metrics:  make(map[string]Metric),
		Errors:   e,
		Interval: interval,
		Hostname: servername,
	}

	prometheus.MustRegister(e)

	if debug {
		log.SetLevel(log.DebugLevel)
	}

	log.Debugf("%+v", metrics)

	go collect(metrics)

	http.Handle("/metrics", promhttp.Handler())
	log.Info("beginning to serve")
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), nil))
}
