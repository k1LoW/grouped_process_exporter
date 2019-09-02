package cgroup

import (
	"os"
	"testing"

	"github.com/k1LoW/grouped_process_exporter/grouped_proc"
	"github.com/k1LoW/grouped_process_exporter/metric"
)

const (
	testProcPath   = "../../testdata/proc"
	testCgroupPath = "../../testdata/sys/fs/cgroup"
)

func TestCollect(t *testing.T) {
	cgroup := testCgroup()
	gpMap := make(map[string]*grouped_proc.GroupedProc)
	enabled := make(map[metric.MetricKey]bool)

	enabled[metric.ProcIO] = true
	enabled[metric.ProcStat] = true
	err := cgroup.Collect(gpMap, enabled)
	if err != nil {
		t.Fatalf("%v", err)
	}

	if len(gpMap) != 2 {
		t.Errorf("want %d, got %d", 2, len(gpMap))
	}
	if _, ok := gpMap["/system.slice/nginx.service"]; !ok {
		t.Errorf("want %s, got none", "/system.slice/nginx.service")
	}
	if _, ok := gpMap["/system.slice/mysql.service"]; !ok {
		t.Errorf("want %s, got none", "/system.slice/mysql.service")
	}
}

func TestCollectWithNormalize(t *testing.T) {
	cgroup := testCgroup()
	err := cgroup.SetNormalizeRegexp("/(system.slice).+")
	if err != nil {
		t.Fatalf("%v", err)
	}
	gpMap := make(map[string]*grouped_proc.GroupedProc)
	enabled := make(map[metric.MetricKey]bool)

	enabled[metric.ProcIO] = true
	enabled[metric.ProcStat] = true
	err = cgroup.Collect(gpMap, enabled)
	if err != nil {
		t.Fatalf("%v", err)
	}

	if len(gpMap) != 1 {
		t.Errorf("want %d, got %d", 1, len(gpMap))
	}
	if _, ok := gpMap["system.slice"]; !ok {
		t.Errorf("want %s, got none", "system.slice")
	}
}

func testCgroup() *Cgroup {
	os.Setenv("GROUPED_PROCESS_PROC_MOUNT_POINT", testProcPath)
	return NewCgroup(testCgroupPath)
}
