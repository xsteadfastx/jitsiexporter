// nolint: goerr113
package jitsiexporter_test

import (
	"errors"
	"fmt"
	"testing"

	"git.xsfx.dev/prometheus/jitsiexporter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestUpdate(t *testing.T) {
	assert := assert.New(t)

	s := make(map[string]interface{})
	s["foo"] = "foo"
	s["bar"] = 1
	s["zonk"] = float64(1)
	mockStater := &jitsiexporter.MockStater{}
	mockStater.On("Now", "http://foo.tld").Return(s, nil)

	m := &jitsiexporter.Metrics{
		URL:     "http://foo.tld",
		Metrics: make(map[string]jitsiexporter.Metric),
		Stater:  mockStater,
	}

	err := m.Update()
	assert.Empty(err)

	assert.Equal(float64(1), testutil.ToFloat64(m.Metrics["jitsi_zonk"].Gauge))
	assert.Equal(jitsiexporter.Metric{Name: "", Gauge: prometheus.Gauge(nil)}, m.Metrics["jitsi_foo"])
	assert.Equal(jitsiexporter.Metric{Name: "", Gauge: prometheus.Gauge(nil)}, m.Metrics["jitsi_bar"])
	assert.Equal(len(m.Metrics), 1)
}

func TestUpdateOnError(t *testing.T) {
	assert := assert.New(t)

	mockStater := &jitsiexporter.MockStater{}
	mockStater.On("Now", "http://foo.tld").Return(nil, errors.New("something went foo"))

	e := prometheus.NewCounter(prometheus.CounterOpts{Name: "jitsi_fetch_errors"})
	metricsMap := make(map[string]jitsiexporter.Metric)
	metricsMap["testmetric"] = jitsiexporter.Metric{
		Name:  "testmetric",
		Gauge: prometheus.NewGauge(prometheus.GaugeOpts{Name: "jitsi_testmetric"}),
	}
	m := &jitsiexporter.Metrics{
		URL:     "http://foo.tld",
		Metrics: metricsMap,
		Stater:  mockStater,
		Errors:  e,
	}

	assert.Equal(1, len(m.Metrics))

	err := m.Update()
	fmt.Println(err)
	assert.NotEmpty(err)

	assert.Equal(0, len(m.Metrics))

	assert.Equal(float64(1), testutil.ToFloat64(e))
}
