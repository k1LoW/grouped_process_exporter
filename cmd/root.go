/*
Copyright © 2019 Ken'ichiro Oyama <k1lowxb@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/k1LoW/grouped_process_exporter/collector"
	"github.com/k1LoW/grouped_process_exporter/grouper"
	"github.com/k1LoW/grouped_process_exporter/grouper/cgroup"
	"github.com/k1LoW/grouped_process_exporter/grouper/proc_status_name"
	"github.com/k1LoW/grouped_process_exporter/metric"
	"github.com/k1LoW/grouped_process_exporter/version"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	address              string
	endpoint             string
	groupType            string
	nReStr               string
	eReStr               string
	collectStat          bool
	collectIO            bool
	collectStatus        bool
	subsystems           []string
	enableMetricDescName string

	format string
	level  string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "grouped_process_exporter",
	Short: "Exporter for grouped process",
	Long:  `Exporter for grouped process.`,
	Args: func(cmd *cobra.Command, args []string) error {
		versionVal, err := cmd.Flags().GetBool("version")
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		if versionVal {
			fmt.Println(version.Version)
			os.Exit(0)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		lvl, err := logrus.ParseLevel(level)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		} else {
			logrus.SetLevel(lvl)
		}
		if format == "json" {
			logrus.SetFormatter(&logrus.JSONFormatter{})
		}

		status, err := runRoot(args, address, endpoint, groupType, nReStr, eReStr, collectStat, collectIO)
		if err != nil {
			logrus.Fatalln(err)
		}
		logrus.Infoln("Stopped grouped_process_exporter")
		os.Exit(status)
	},
}

func runRoot(args []string, address, endpoint, groupType, nReStr, eReStr string, collectStat, collectIO bool) (int, error) {
	var g grouper.Grouper
	switch groupType {
	case "cgroup":
		logrus.Infoln("Select cgroup grouper")
		fsPath := "/sys/fs/cgroup"
		g = cgroup.NewCgroup(fsPath, subsystems)
	case "proc_status_name", "name":
		logrus.Infoln("Select proc_status_name grouper")
		g = proc_status_name.NewProcStatusName()
	default:
		return 1, errors.New("invalid grouping type")
	}
	err := g.SetNormalizeRegexp(nReStr)
	if err != nil {
		return 1, err
	}
	err = g.SetExcludeRegexp(eReStr)
	if err != nil {
		return 1, err
	}

	collector, err := collector.NewGroupedProcCollector(g)
	if err != nil {
		return 1, err
	}
	if collectStat {
		collector.EnableMetric(metric.ProcStat)
		logrus.Infoln("Enable collecting /proc/[PID]/stat.")
	}
	if collectIO {
		collector.EnableMetric(metric.ProcIO)
		logrus.Infoln("Enable collecting /proc/[PID]/io.")
	}
	if collectStatus {
		collector.EnableMetric(metric.ProcStatus)
		logrus.Infoln("Enable collecting /proc/[PID]/status.")
	}
	if err := collector.SetEnableMetricDescNameRegexp(enableMetricDescName); err != nil {
		return 1, err
	}

	r := prometheus.NewRegistry()
	r.MustRegister(collectors.NewBuildInfoCollector())
	if err := r.Register(collector); err != nil {
		return 1, fmt.Errorf("couldn't register grouped_process_collector: %s", err)
	}

	srv := &http.Server{
		Addr: address,
	}

	handler := promhttp.HandlerFor(
		prometheus.Gatherers{r},
		promhttp.HandlerOpts{
			ErrorLog:            log.New(logrus.StandardLogger().Writer(), "", 0),
			ErrorHandling:       promhttp.ContinueOnError,
			MaxRequestsInFlight: 10,
			Registry:            r,
		},
	)

	http.Handle(endpoint, handler)
	logrus.Infoln("Starting grouped_process_exporter", version.Version)
	logrus.Infoln(fmt.Sprintf("Listening on %s%s", address, endpoint))
	if err := srv.ListenAndServe(); err != nil {
		return 1, err
	}
	return 0, nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatalln(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&address, "telemetry.address", "", ":9644", "Address on which to expose metrics.")
	rootCmd.Flags().StringVarP(&endpoint, "telemetry.endpoint", "", "/metrics", "Path under which to expose metrics.")
	rootCmd.Flags().StringVarP(&groupType, "group.type", "", "cgroup", "Grouping type.")
	rootCmd.Flags().StringVarP(&nReStr, "group.normalize", "", "", "Regexp for normalize group names. Exporter use regexp match result $1 as group name.")
	rootCmd.Flags().StringVarP(&eReStr, "group.exclude", "", "", "Regexp for exclude group names. Exporter exclude group using regexp match before group name normalization")
	rootCmd.Flags().BoolVarP(&collectStat, "collector.stat", "", false, "Enable collecting /proc/[PID]/stat.")
	rootCmd.Flags().BoolVarP(&collectIO, "collector.io", "", false, "Enable collecting /proc/[PID]/io.")
	rootCmd.Flags().BoolVarP(&collectStatus, "collector.status", "", false, "Enable collecting /proc/[PID]/status.")
	rootCmd.Flags().StringArrayVarP(&subsystems, "cgroup.subsystem", "", []string{}, fmt.Sprintf("Cgroup subsystem to scan. (default %s)", cgroup.DefaultSubsystems))
	rootCmd.Flags().StringVarP(&enableMetricDescName, "metric.desc", "", ".+", "Regexp for enable metric descriptor.")

	rootCmd.Flags().StringVarP(&level, "log.level", "", logrus.New().Level.String(), "Only log messages with the given severity or above. Valid levels: [debug, info, warn, error, fatal]")
	rootCmd.Flags().StringVarP(&format, "log.format", "", "text", `Set the log format. Valid formats: [text, json]`)

	rootCmd.Flags().BoolP("version", "v", false, "print the version")
}
