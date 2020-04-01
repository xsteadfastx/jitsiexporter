package jitsiexporter

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestUpdate(t *testing.T) {
	assert := assert.New(t)

	s := make(map[string]interface{})
	s["foo"] = "foo"
	s["bar"] = 1           // nolint:gomnd
	s["zonk"] = float64(1) // nolint:gomnd
	mockStater := &MockStater{}
	mockStater.On("Now", "http://foo.tld").Return(s, nil)

	m := &Metrics{
		URL:     "http://foo.tld",
		Metrics: make(map[string]Metric),
		Stater:  mockStater,
	}

	err := m.Update()
	assert.Empty(err)

	assert.Equal(testutil.ToFloat64(m.Metrics["jitsi_zonk"].Gauge), float64(1))
	assert.Equal(m.Metrics["jitsi_foo"], Metric{Name: "", Gauge: prometheus.Gauge(nil)})
	assert.Equal(m.Metrics["jitsi_bar"], Metric{Name: "", Gauge: prometheus.Gauge(nil)})
	assert.Equal(len(m.Metrics), 1)
}
