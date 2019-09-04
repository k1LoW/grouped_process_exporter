package proc_status_name

import (
	"os"
	"testing"

	"github.com/k1LoW/grouped_process_exporter/grouped_proc"
	"github.com/k1LoW/grouped_process_exporter/metric"
)

const (
	testProcPath = "../../testdata/proc"
)

func TestCollect(t *testing.T) {
	procStatusName := testProcStatusName()
	gprocs := grouped_proc.NewGroupedProcs()
	enabled := make(map[metric.MetricKey]bool)

	enabled[metric.ProcIO] = true
	enabled[metric.ProcStat] = true
	err := procStatusName.Collect(gprocs, enabled)
	if err != nil {
		t.Fatalf("%v", err)
	}

	if gprocs.Length() != 2 {
		t.Errorf("want %d, got %d", 2, gprocs.Length())
	}
	if _, ok := gprocs.Load("nginx"); !ok {
		t.Errorf("want %s, got none", "nginx")
	}
	if _, ok := gprocs.Load("mysqld"); !ok {
		t.Errorf("want %s, got none", "mysqld")
	}
}

func TestCollectWithNormalize(t *testing.T) {
	procStatusName := testProcStatusName()
	err := procStatusName.SetNormalizeRegexp("(mysql).+")
	if err != nil {
		t.Fatalf("%v", err)
	}
	gprocs := grouped_proc.NewGroupedProcs()
	enabled := make(map[metric.MetricKey]bool)

	enabled[metric.ProcIO] = true
	enabled[metric.ProcStat] = true
	err = procStatusName.Collect(gprocs, enabled)
	if err != nil {
		t.Fatalf("%v", err)
	}

	if gprocs.Length() != 2 {
		t.Errorf("want %d, got %d", 2, gprocs.Length())
	}
	if _, ok := gprocs.Load("nginx"); !ok {
		t.Errorf("want %s, got none", "nginx")
	}
	if _, ok := gprocs.Load("mysql"); !ok {
		t.Errorf("want %s, got none", "mysql")
	}
}

func testProcStatusName() *ProcStatusName {
	os.Setenv("GROUPED_PROCESS_PROC_MOUNT_POINT", testProcPath)
	return NewProcStatusName()
}
