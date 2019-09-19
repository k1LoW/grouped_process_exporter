## [v0.4.0](https://github.com/k1LoW/grouped_process_exporter/compare/v0.3.0...v0.4.0) (2019-09-19)

* Add `--group.exclude` option for exclude group using regexp [#13](https://github.com/k1LoW/grouped_process_exporter/pull/13) ([k1LoW](https://github.com/k1LoW))

## [v0.3.0](https://github.com/k1LoW/grouped_process_exporter/compare/v0.2.0...v0.3.0) (2019-09-17)

* Use prometheus/procfs v0.0.5 [#12](https://github.com/k1LoW/grouped_process_exporter/pull/12) ([k1LoW](https://github.com/k1LoW))

## [v0.2.0](https://github.com/k1LoW/grouped_process_exporter/compare/v0.1.5...v0.2.0) (2019-09-05)

* Change desc `grouped_process_procs` to `grouped_process_num_procs` [#11](https://github.com/k1LoW/grouped_process_exporter/pull/11) ([k1LoW](https://github.com/k1LoW))

## [v0.1.5](https://github.com/k1LoW/grouped_process_exporter/compare/v0.1.4...v0.1.5) (2019-09-05)

* Fix proc_io_syscw calc bug [#10](https://github.com/k1LoW/grouped_process_exporter/pull/10) ([k1LoW](https://github.com/k1LoW))
* Fix proc_stat_numthreads/proc_stat_vsize/proc_stat_rss calc bug [#9](https://github.com/k1LoW/grouped_process_exporter/pull/9) ([k1LoW](https://github.com/k1LoW))
* Fix collector concurrent bug [#8](https://github.com/k1LoW/grouped_process_exporter/pull/8) ([k1LoW](https://github.com/k1LoW))
* Fix metric/ concurrent map iteration and map write [#7](https://github.com/k1LoW/grouped_process_exporter/pull/7) ([k1LoW](https://github.com/k1LoW))
* Fix concurrent map (map[string]*grouped_proc.GroupedProc) writes [#6](https://github.com/k1LoW/grouped_process_exporter/pull/6) ([k1LoW](https://github.com/k1LoW))
* Add github.com/prometheus/common/log [#5](https://github.com/k1LoW/grouped_process_exporter/pull/5) ([k1LoW](https://github.com/k1LoW))
* Add test [#4](https://github.com/k1LoW/grouped_process_exporter/pull/4) ([k1LoW](https://github.com/k1LoW))
* Add `--version` option for print version [#3](https://github.com/k1LoW/grouped_process_exporter/pull/3) ([k1LoW](https://github.com/k1LoW))
* Add grouped /proc/[PID]/stat metrics [#2](https://github.com/k1LoW/grouped_process_exporter/pull/2) ([k1LoW](https://github.com/k1LoW))
* Fix counting architecture [#1](https://github.com/k1LoW/grouped_process_exporter/pull/1) ([k1LoW](https://github.com/k1LoW))

## [v0.1.4](https://github.com/k1LoW/grouped_process_exporter/compare/v0.1.3...v0.1.4) (2019-09-05)

* Fix proc_stat_numthreads/proc_stat_vsize/proc_stat_rss calc bug [#9](https://github.com/k1LoW/grouped_process_exporter/pull/9) ([k1LoW](https://github.com/k1LoW))

## [v0.1.3](https://github.com/k1LoW/grouped_process_exporter/compare/v0.1.2...v0.1.3) (2019-09-04)

* Fix collector concurrent bug [#8](https://github.com/k1LoW/grouped_process_exporter/pull/8) ([k1LoW](https://github.com/k1LoW))

## [v0.1.2](https://github.com/k1LoW/grouped_process_exporter/compare/v0.1.1...v0.1.2) (2019-09-04)

* Fix metric/ concurrent map iteration and map write [#7](https://github.com/k1LoW/grouped_process_exporter/pull/7) ([k1LoW](https://github.com/k1LoW))

## [v0.1.1](https://github.com/k1LoW/grouped_process_exporter/compare/v0.1.0...v0.1.1) (2019-09-04)

* Fix concurrent map (map[string]*grouped_proc.GroupedProc) writes [#6](https://github.com/k1LoW/grouped_process_exporter/pull/6) ([k1LoW](https://github.com/k1LoW))

## [v0.1.0](https://github.com/k1LoW/grouped_process_exporter/compare/0b50674837e9...v0.1.0) (2019-09-03)

* Add github.com/prometheus/common/log [#5](https://github.com/k1LoW/grouped_process_exporter/pull/5) ([k1LoW](https://github.com/k1LoW))
* Add test [#4](https://github.com/k1LoW/grouped_process_exporter/pull/4) ([k1LoW](https://github.com/k1LoW))
* Add `--version` option for print version [#3](https://github.com/k1LoW/grouped_process_exporter/pull/3) ([k1LoW](https://github.com/k1LoW))
* Add grouped /proc/[PID]/stat metrics [#2](https://github.com/k1LoW/grouped_process_exporter/pull/2) ([k1LoW](https://github.com/k1LoW))
* Fix counting architecture [#1](https://github.com/k1LoW/grouped_process_exporter/pull/1) ([k1LoW](https://github.com/k1LoW))
