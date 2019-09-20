package cgroup

import (
	"os"
	"testing"

	"github.com/k1LoW/grouped_process_exporter/grouped_proc"
	"github.com/k1LoW/grouped_process_exporter/metric"
	"golang.org/x/sync/semaphore"
)

const (
	testProcPath   = "../../testdata/proc"
	testCgroupPath = "../../testdata/sys/fs/cgroup"
)

func TestCollect(t *testing.T) {
	cgroup := testCgroup()
	gprocs := grouped_proc.NewGroupedProcs()
	enabled := make(map[metric.MetricKey]bool)
	sem := semaphore.NewWeighted(5)

	enabled[metric.ProcIO] = true
	enabled[metric.ProcStat] = true
	err := cgroup.Collect(gprocs, enabled, sem)
	if err != nil {
		t.Fatalf("%v", err)
	}

	if gprocs.Length() != 2 {
		t.Errorf("want %d, got %d", 2, gprocs.Length())
	}
	if _, ok := gprocs.Load("/system.slice/nginx.service"); !ok {
		t.Errorf("want %s, got none", "/system.slice/nginx.service")
	}
	if _, ok := gprocs.Load("/system.slice/mysql.service"); !ok {
		t.Errorf("want %s, got none", "/system.slice/mysql.service")
	}
}

func TestCollectWithNormalize(t *testing.T) {
	cgroup := testCgroup()
	err := cgroup.SetNormalizeRegexp("/(system.slice).+")
	if err != nil {
		t.Fatalf("%v", err)
	}
	gprocs := grouped_proc.NewGroupedProcs()
	enabled := make(map[metric.MetricKey]bool)
	sem := semaphore.NewWeighted(5)

	enabled[metric.ProcIO] = true
	enabled[metric.ProcStat] = true
	err = cgroup.Collect(gprocs, enabled, sem)
	if err != nil {
		t.Fatalf("%v", err)
	}

	if gprocs.Length() != 1 {
		t.Errorf("want %d, got %d", 1, gprocs.Length())
	}
	if _, ok := gprocs.Load("system.slice"); !ok {
		t.Errorf("want %s, got none", "system.slice")
	}
}

func TestCollectWithExclude(t *testing.T) {
	cgroup := testCgroup()
	err := cgroup.SetExcludeRegexp("my.+")
	if err != nil {
		t.Fatalf("%v", err)
	}
	gprocs := grouped_proc.NewGroupedProcs()
	enabled := make(map[metric.MetricKey]bool)
	sem := semaphore.NewWeighted(5)

	enabled[metric.ProcIO] = true
	enabled[metric.ProcStat] = true
	err = cgroup.Collect(gprocs, enabled, sem)
	if err != nil {
		t.Fatalf("%v", err)
	}

	if gprocs.Length() != 1 {
		t.Errorf("want %d, got %d", 1, gprocs.Length())
	}
	if _, ok := gprocs.Load("/system.slice/nginx.service"); !ok {
		t.Errorf("want %s, got none", "/system.slice/nginx.service")
	}
}

func testCgroup() *Cgroup {
	os.Setenv("GROUPED_PROCESS_PROC_MOUNT_POINT", testProcPath)
	return NewCgroup(testCgroupPath)
}
