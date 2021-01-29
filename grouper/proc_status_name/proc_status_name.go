package proc_status_name

import (
	"context"
	"errors"
	"os"
	"regexp"
	"sync"

	"github.com/k1LoW/grouped_process_exporter/grouped_proc"
	"github.com/k1LoW/grouped_process_exporter/metric"
	"github.com/prometheus/procfs"
	"golang.org/x/sync/semaphore"
)

type ProcStatusName struct {
	nRe            *regexp.Regexp // normalize regexp
	eRe            *regexp.Regexp // exclude regexp
	procMountPoint string
}

func (g *ProcStatusName) Name() string {
	return "proc_status_name"
}

func (g *ProcStatusName) Collect(gprocs *grouped_proc.GroupedProcs, enabled map[metric.MetricKey]bool, sem *semaphore.Weighted) error {
	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_ = sem.Acquire(ctx, 1)
	fs, err := procfs.NewFS(g.procMountPoint)
	if err != nil {
		sem.Release(1)
		return err
	}
	procs, err := fs.AllProcs()
	if err != nil {
		sem.Release(1)
		return err
	}
	sem.Release(1)

	for _, proc := range procs {
		_ = sem.Acquire(ctx, 1)
		status, err := proc.NewStatus()
		sem.Release(1)
		if err != nil {
			// TODO: Log
			continue
		}

		// collect process only
		if status.PID != status.TGID {
			continue
		}

		pid := proc.PID
		name := status.Name
		if g.eRe != nil {
			if g.eRe.MatchString(name) {
				continue
			}
		}
		if g.nRe != nil {
			matches := g.nRe.FindStringSubmatch(name)
			if len(matches) > 1 {
				name = matches[1]
			}
		}
		var (
			gproc *grouped_proc.GroupedProc
			ok    bool
		)
		gproc, ok = gprocs.Load(name)
		if !ok {
			gproc = grouped_proc.NewGroupedProc(enabled)
			gprocs.Store(name, gproc)
		}
		if err := gproc.Collect(name); err != nil {
			return err
		}
		gproc.Exists = true
		wg.Add(1)
		_ = sem.Acquire(ctx, gproc.RequiredWeight)
		go func(wg *sync.WaitGroup, pid int, gproc *grouped_proc.GroupedProc) {
			_ = gproc.AppendProcAndCollect(pid)
			sem.Release(gproc.RequiredWeight)
			wg.Done()
		}(wg, pid, gproc)
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

func (g *ProcStatusName) SetExcludeRegexp(eReStr string) error {
	if eReStr == "" {
		return nil
	}
	eRe, err := regexp.Compile(eReStr)
	if err != nil {
		return err
	}
	g.eRe = eRe
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
