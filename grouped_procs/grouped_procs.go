package grouped_procs

import "github.com/prometheus/procfs"

type GroupedProcs struct {
	Path string
	Pids []int
	IO   procfs.ProcIO
}

func (g *GroupedProcs) AppendPid(pid int) error {
	proc, err := procfs.NewProc(pid)
	if err != nil {
		return err
	}

	// cgroup.procs (PIDs)
	g.Pids = append(g.Pids, pid)

	// /proc/[PID]/io
	io, err := proc.IO()
	if err != nil {
		return err
	}
	g.IO.RChar = g.IO.RChar + io.RChar
	g.IO.WChar = g.IO.WChar + io.WChar
	g.IO.SyscR = g.IO.SyscR + io.SyscR
	g.IO.SyscW = g.IO.SyscW + io.SyscW
	g.IO.ReadBytes = g.IO.ReadBytes + io.ReadBytes
	g.IO.WriteBytes = g.IO.WriteBytes + io.WriteBytes
	g.IO.CancelledWriteBytes = g.IO.CancelledWriteBytes + io.CancelledWriteBytes
	return nil
}
