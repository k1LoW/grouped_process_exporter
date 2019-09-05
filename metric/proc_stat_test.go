package metric

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/procfs"
)

func TestProcStatMetric(t *testing.T) {
	m := NewProcStatMetric()
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

func TestProcStatCollectFromProc(t *testing.T) {
	pids := []int{1357, 1380, 1383, 1384, 1388}

	m := NewProcStatMetric()
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
		case strings.Contains(desc, "grouped_process_stat_minflt_total"):
			got = d.GetCounter().GetValue()
			want = float64(46 + 787 + 331 + 635 + 756)
		case strings.Contains(desc, "grouped_process_stat_cminflt_total"):
			got = d.GetCounter().GetValue()
			want = float64(0 + 0 + 0 + 0 + 0)
		case strings.Contains(desc, "grouped_process_stat_majflt_total"):
			got = d.GetCounter().GetValue()
			want = float64(0 + 1 + 0 + 9 + 0)
		case strings.Contains(desc, "grouped_process_stat_cmajflt_total"):
			got = d.GetCounter().GetValue()
			want = float64(0 + 0 + 0 + 0 + 0)
		case strings.Contains(desc, "grouped_process_stat_utime_total"):
			got = d.GetCounter().GetValue()
			want = float64(0 + 17 + 12 + 6 + 16)
		case strings.Contains(desc, "grouped_process_stat_stime_total"):
			got = d.GetCounter().GetValue()
			want = float64(0 + 312 + 97 + 192 + 304)
		case strings.Contains(desc, "grouped_process_stat_cutime_total"):
			got = d.GetCounter().GetValue()
			want = float64(0 + 0 + 0 + 0 + 0)
		case strings.Contains(desc, "grouped_process_stat_cstime_total"):
			got = d.GetCounter().GetValue()
			want = float64(0 + 0 + 0 + 0 + 0)
		case strings.Contains(desc, "grouped_process_stat_numthreads"):
			got = d.GetGauge().GetValue()
			want = float64(1 + 1 + 1 + 1 + 1)
		case strings.Contains(desc, "grouped_process_stat_vsize_bytes"):
			got = d.GetGauge().GetValue()
			want = float64(128118784 + 129372160 + 128507904 + 128507904 + 128507904)
		case strings.Contains(desc, "grouped_process_stat_rss"):
			got = d.GetGauge().GetValue()
			want = float64(333 + 1274 + 982 + 982 + 998)
		case strings.Contains(desc, "grouped_process_stat_clk_tck"):
			got = d.GetGauge().GetValue()
			want = ClkTck()
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
