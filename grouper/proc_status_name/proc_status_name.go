package proc_status_name

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/k1LoW/grouped_process_exporter/grouped_proc"
	"github.com/k1LoW/grouped_process_exporter/metric"
	"github.com/prometheus/procfs"
)

type ProcStatusName struct {
	nRe            *regexp.Regexp
	procMountPoint string
}

func (g *ProcStatusName) Name() string {
	return "proc_status_name"
}

func (g *ProcStatusName) Collect(gpMap map[string]*grouped_proc.GroupedProc, enabled map[metric.MetricKey]bool) error {
	wg := &sync.WaitGroup{}
	fs, err := procfs.NewFS(g.procMountPoint)
	if err != nil {
		return err
	}
	procs, err := fs.AllProcs()
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", procs)

	for _, proc := range procs {
		status, err := proc.NewStatus()
		if err != nil {
			// TODO: Log
			continue
		}
		pid := proc.PID

		// collect process only
		b, err := ioutil.ReadFile(filepath.Join(g.procMountPoint, strconv.Itoa(pid), "status"))
		if err != nil {
			continue
		}
		if !strings.Contains(string(b), fmt.Sprintf("Tgid:\t%d", pid)) {
			continue
		}
		name := status.Name
		if g.nRe != nil {
			matches := g.nRe.FindStringSubmatch(name)
			if len(matches) > 1 {
				name = matches[1]
			}
		}
		_, ok := gpMap[name]
		if !ok {
			gpMap[name] = grouped_proc.NewGroupedProc(enabled)
		}
		gpMap[name].Exists = true
		wg.Add(1)
		go func(wg *sync.WaitGroup, pid int, g *grouped_proc.GroupedProc) {
			_ = g.AppendProcAndCollect(pid)
			wg.Done()
		}(wg, pid, gpMap[name])
	}
	wg.Wait()
	return nil
}

func (g *ProcStatusName) SetNormalizeRegexp(nReStr string) error {
	if nReStr == "" {
		return nil
	}
	nRe, err := regexp.Compile(nReStr)
	if err != nil {
		return err
	}
	if nRe.NumSubexp() != 1 {
		return errors.New("number of parenthesized subexpressions in this regexp should be 1")
	}
	g.nRe = nRe
	return nil
}

// NewProcStatusName
func NewProcStatusName() *ProcStatusName {
	procMountPoint := os.Getenv("GROUPED_PROCESS_PROC_MOUNT_POINT")
	if procMountPoint == "" {
		procMountPoint = procfs.DefaultMountPoint
	}
	return &ProcStatusName{
		procMountPoint: procMountPoint,
	}
}
