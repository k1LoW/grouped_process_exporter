package metric

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/procfs"
)

func TestProcStatusMetric(t *testing.T) {
	m := NewProcStatusMetric()
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

func TestProcStatusCollectFromProc(t *testing.T) {
	pids := []int{1357, 1380, 1383, 1384, 1388}

	m := NewProcStatusMetric()
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
		case strings.Contains(desc, "grouped_process_status_VmPeak_bytes_total"):
			got = d.GetGauge().GetValue()
			want = float64((125116 + 126424 + 125760 + 126408 + 126420) * 1024)
		case strings.Contains(desc, "grouped_process_status_VmSize_bytes_total"):
			got = d.GetGauge().GetValue()
			want = float64((125116 + 126340 + 125496 + 125496 + 125496) * 1024)
		case strings.Contains(desc, "grouped_process_status_VmLck_bytes_total"):
			got = d.GetGauge().GetValue()
			want = float64((0 + 0 + 0 + 0 + 0) * 1024)
		case strings.Contains(desc, "grouped_process_status_VmPin_bytes_total"):
			got = d.GetGauge().GetValue()
			want = float64((0 + 0 + 0 + 0 + 0) * 1024)
		case strings.Contains(desc, "grouped_process_status_VmHWM_bytes_total"):
			got = d.GetGauge().GetValue()
			want = float64((1516 + 5900 + 4124 + 4872 + 5012) * 1024)
		case strings.Contains(desc, "grouped_process_status_VmRSS_bytes_total"):
			got = d.GetGauge().GetValue()
			want = float64((1332 + 5096 + 3928 + 3928 + 3992) * 1024)
		case strings.Contains(desc, "grouped_process_status_RssAnon_bytes_total"):
			got = d.GetGauge().GetValue()
			want = float64((0 + 0 + 0 + 0 + 0) * 1024)
		case strings.Contains(desc, "grouped_process_status_RssFile_bytes_total"):
			got = d.GetGauge().GetValue()
			want = float64((0 + 0 + 0 + 0 + 0) * 1024)
		case strings.Contains(desc, "grouped_process_status_RssShmem_bytes_total"):
			got = d.GetGauge().GetValue()
			want = float64((0 + 0 + 0 + 0 + 0) * 1024)
		case strings.Contains(desc, "grouped_process_status_VmData_bytes_total"):
			got = d.GetGauge().GetValue()
			want = float64((1088 + 2312 + 1468 + 1468 + 1468) * 1024)
		case strings.Contains(desc, "grouped_process_status_VmStk_bytes_total"):
			got = d.GetGauge().GetValue()
			want = float64((132 + 132 + 132 + 132 + 132) * 1024)
		case strings.Contains(desc, "grouped_process_status_VmExe_bytes_total"):
			got = d.GetGauge().GetValue()
			want = float64((1080 + 1080 + 1080 + 1080 + 1080) * 1024)
		case strings.Contains(desc, "grouped_process_status_VmLib_bytes_total"):
			got = d.GetGauge().GetValue()
			want = float64((44076 + 44076 + 44076 + 44076 + 44076) * 1024)
		case strings.Contains(desc, "grouped_process_status_VmPTE_bytes_total"):
			got = d.GetGauge().GetValue()
			want = float64((180 + 180 + 180 + 180 + 180) * 1024)
		case strings.Contains(desc, "grouped_process_status_VmPMD_bytes_total"):
			got = d.GetGauge().GetValue()
			want = float64((12 + 12 + 12 + 12 + 12) * 1024)
		case strings.Contains(desc, "grouped_process_status_VmSwap_bytes_total"):
			got = d.GetGauge().GetValue()
			want = float64((100 + 72 + 76 + 76 + 76) * 1024)
		case strings.Contains(desc, "grouped_process_status_HugetlbPages_bytes_total"):
			got = d.GetGauge().GetValue()
			want = float64((0 + 0 + 0 + 0 + 0) * 1024)
		case strings.Contains(desc, "grouped_process_status_VoluntaryCtxtSwitches_total"):
			got = d.GetCounter().GetValue()
			want = float64(1 + 51333 + 35604 + 14718 + 53736)
		case strings.Contains(desc, "grouped_process_status_NonVoluntaryCtxtSwitches_total"):
			got = d.GetCounter().GetValue()
			want = float64(4 + 626 + 102 + 529 + 730)
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
