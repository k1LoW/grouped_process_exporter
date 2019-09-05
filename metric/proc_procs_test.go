package metric

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/procfs"
)

const (
	testProcPath = "../testdata/proc"
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

func TestProcProcsCollectFromProc(t *testing.T) {
	pids := []int{1357, 1380, 1383, 1384, 1388}

	m := NewProcProcsMetric()
	descs := m.Describe()
	ch := make(chan prometheus.Metric, 100)

	fs, err := procfs.NewFS(testProcPath)
	if err != nil {
		t.Fatalf("%v", err)
	}
	for _, pid := range pids {
		proc, err := fs.Proc(pid)
		if err != nil {
			t.Fatalf("%v", err)
		}
		err = m.CollectFromProc(proc)
		if err != nil {
			t.Fatalf("%v", err)
		}
	}
	err = m.PushCollected(ch, descs, "cgroup", "test")
	if err != nil {
		t.Fatalf("%v", err)
	}
	for m := range ch {
		d := &dto.Metric{}
		err := m.Write(d)
		if err != nil {
			t.Fatalf("%v", err)
		}
		desc := m.Desc().String()
		var (
			got  float64
			want float64
		)
		switch {
		case strings.Contains(desc, "grouped_process_num_procs"):
			got = d.GetGauge().GetValue()
			want = float64(len(pids))
		default:
			t.Fatalf("unknowns desc :%s", desc)
		}
		if want != got {
			t.Errorf("%s: want %f, got %f", desc, want, got)
		}
		if len(ch) == 0 {
			close(ch)
		}
	}
}
