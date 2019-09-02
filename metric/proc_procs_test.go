package metric

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func TestProcProcsMetric(t *testing.T) {
	m := NewProcProcsMetric()
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
