# Exporter for grouped process [![GitHub release](https://img.shields.io/github/release/k1LoW/grouped_process_exporter.svg)](https://github.com/k1LoW/grouped_process_exporter/releases) [![GitHub Actions](https://action-badges.now.sh/k1LoW/grouped_process_exporter)](https://github.com/k1LoW/grouped_process_exporter/actions) [![codecov](https://codecov.io/gh/k1LoW/grouped_process_exporter/branch/master/graph/badge.svg)](https://codecov.io/gh/k1LoW/grouped_process_exporter)

## Supported grouping

- control group v1 ( `cgroup`, default )
- /proc/[PID]/status.Name: ( `proc_status_name` / `name` )

### :unlock: Advanced grouping

#### `--group.normalize`

Exporter normalize group names using regexp `--group.normalize` option.

For example, by setting `--group.normalize='^.+(tcpdp).+$'`, exporter normaliz the group names `/path/to/tcpdp-eth0` and `/path/to/tcpdp-eth1` to `tcpdp`.

#### `--group.exclude`

Exporter exclude group using regexp `--group.exclude` option.

For example, by setting `--group.exclude='user.\slice'`, exporter exclude the group names `/user.slice` , `/user.slice/user-10503.slice`.

> Exporter exclude group before group name normalization.

## Avairable Metrics

### Amount of grouped procs ( default on )

| Name | Type | Description |
| --- | --- | --- |
| grouped_process_num_procs | Gauge | Number of processes in the group |

### Grouped /proc/[PID]/stat ( `--collector.stat` )

| Name | Type | Description |
| --- | --- | --- |
| grouped_process_stat_minflt_total | Counter | Total number of grouped /proc/[PID]/stat.minflt |
| grouped_process_stat_cminflt_total | Counter | Total number of grouped /proc/[PID]/stat.rchar |
| grouped_process_stat_majflt_total | Counter | Total number of grouped /proc/[PID]/stat.majflt |
| grouped_process_stat_cmajflt_total | Counter | Total number of grouped /proc/[PID]/stat.cmajflt |
| grouped_process_stat_utime_total | Counter | Total number of grouped /proc/[PID]/stat.utime |
| grouped_process_stat_stime_total | Counter | Total number of grouped /proc/[PID]/stat.stime |
| grouped_process_stat_cutime_total | Counter | Total number of grouped /proc/[PID]/stat.cutime |
| grouped_process_stat_cstime_total | Counter | Total number of grouped /proc/[PID]/stat.cstime |
| grouped_process_stat_numthreads | Gauge | Grouped /proc/[PID]/stat.numthreads |
| grouped_process_stat_vsize_bytes | Gauge | Grouped /proc/[PID]/stat.vsize |
| grouped_process_stat_rss | Gauge | Grouped /proc/[PID]/stat.rss |
| grouped_process_stat_clk_tck | Gauge | clock ticks (divide by sysconf(_SC_CLK_TCK)) |

### Grouped /proc/[PID]/io ( `--collector.io` )

| Name | Type | Description |
| --- | --- | --- |
| grouped_process_io_r_char_total | Counter | Total number of grouped /proc/[PID]/io.rchar |
| grouped_process_io_w_char_total | Counter | Total number of grouped /proc/[PID]/io.wchar |
| grouped_process_io_sysc_r_total | Counter | Total number of grouped /proc/[PID]/io.syscr |
| grouped_process_io_sysc_w_total | Counter | Total number of grouped /proc/[PID]/io.syscw |
| grouped_process_io_read_bytes_total | Counter | Total number of grouped /proc/[PID]/io.read_bytes |
| grouped_process_io_write_bytes_total | Counter | Total number of grouped /proc/[PID]/io.write_bytes |
| grouped_process_io_cancelled_write_bytes_total | Counter | Total number of grouped /proc/[PID]/io.cancelled_write_bytes |

## Alternatives

- [process-exporter](https://github.com/ncabatoff/process-exporter): Prometheus exporter that mines /proc to report on selected processes
