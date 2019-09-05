package metric

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/procfs"
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

func TestProcIOCollectFromProc(t *testing.T) {
	pids := []int{1357, 1380, 1383, 1384, 1388}

	m := NewProcIOMetric()
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
		case strings.Contains(desc, "grouped_process_io_r_char_total"):
			got = d.GetCounter().GetValue()
			want = float64(0 + 4204878 + 534390 + 2968142 + 4054374)
		case strings.Contains(desc, "grouped_process_io_w_char_total"):
			got = d.GetCounter().GetValue()
			want = float64(5 + 12568999 + 1593288 + 8976934 + 12125523)
		case strings.Contains(desc, "grouped_process_io_sysc_r_total"):
			got = d.GetCounter().GetValue()
			want = float64(0 + 5141 + 657 + 3637 + 4927)
		case strings.Contains(desc, "grouped_process_io_sysc_w_total"):
			got = d.GetCounter().GetValue()
			want = float64(1 + 6955 + 876 + 5234 + 6802)
		case strings.Contains(desc, "grouped_process_io_read_bytes_total"):
			got = d.GetCounter().GetValue()
			want = float64(0 + 40960 + 0 + 794624 + 0)
		case strings.Contains(desc, "grouped_process_io_write_bytes_total"):
			got = d.GetCounter().GetValue()
			want = float64(0 + 196608 + 32768 + 143360 + 212992)
		case strings.Contains(desc, "grouped_process_io_cancelled_write_bytes_total"):
			got = d.GetCounter().GetValue()
			want = float64(0 + 0 + 0 + 0 + 0)
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
