package metric

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func TestProcIOMetric(t *testing.T) {
	m := NewProcIOMetric()
	descs := m.Describe()
	ch := make(chan prometheus.Metric, 100)
	err := m.PushCollected(ch, descs, "cgroup", "test")
	if err != nil {
		t.Fatalf("%v", err)
	}
	if len(descs) != len(ch) {
		t.Errorf("descs:%d != ch:%d", len(descs), len(ch))
	}
}