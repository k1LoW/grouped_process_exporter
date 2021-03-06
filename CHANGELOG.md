## [v0.8.0](https://github.com/k1LoW/grouped_process_exporter/compare/v0.7.1...v0.8.0) (2021-01-29)

* Remove timeout for http server. [#29](https://github.com/k1LoW/grouped_process_exporter/pull/29) ([k1LoW](https://github.com/k1LoW))
* Add default metrics `grouped_process_num_grouped` ( Number of grouped ) [#28](https://github.com/k1LoW/grouped_process_exporter/pull/28) ([k1LoW](https://github.com/k1LoW))
* Add --metric.desc to filter metric descriptor [#27](https://github.com/k1LoW/grouped_process_exporter/pull/27) ([k1LoW](https://github.com/k1LoW))
* Bump up go and pkg version [#26](https://github.com/k1LoW/grouped_process_exporter/pull/26) ([k1LoW](https://github.com/k1LoW))

## [v0.7.1](https://github.com/k1LoW/grouped_process_exporter/compare/v0.7.0...v0.7.1) (2021-01-27)

* Fix invalid memory address or nil pointer dereference when --collector.status is enabled [#25](https://github.com/k1LoW/grouped_process_exporter/pull/25) ([k1LoW](https://github.com/k1LoW))

## [v0.7.0](https://github.com/k1LoW/grouped_process_exporter/compare/v0.6.0...v0.7.0) (2019-12-11)

* Update prometheus/procfs to v0.0.8 [#24](https://github.com/k1LoW/grouped_process_exporter/pull/24) ([k1LoW](https://github.com/k1LoW))
* Add grouped /proc/[PID]/status metrics [#23](https://github.com/k1LoW/grouped_process_exporter/pull/23) ([k1LoW](https://github.com/k1LoW))

## [v0.6.0](https://github.com/k1LoW/grouped_process_exporter/compare/v0.5.1...v0.6.0) (2019-11-14)

* Resolve cgroup subsystem symlinks [#22](https://github.com/k1LoW/grouped_process_exporter/pull/22) ([k1LoW](https://github.com/k1LoW))
* Fix `--cgroup.subsystem` option type [#21](https://github.com/k1LoW/grouped_process_exporter/pull/21) ([k1LoW](https://github.com/k1LoW))
* Add `--cgroup.subsystem` option to set cgroup subsystem to scan [#20](https://github.com/k1LoW/grouped_process_exporter/pull/20) ([k1LoW](https://github.com/k1LoW))
* Update prometheus/procfs to v0.0.6 [#19](https://github.com/k1LoW/grouped_process_exporter/pull/19) ([k1LoW](https://github.com/k1LoW))
* README.md: fix typo [#18](https://github.com/k1LoW/grouped_process_exporter/pull/18) ([perlun](https://github.com/perlun))
* README.md: minor English and grammar improvements [#17](https://github.com/k1LoW/grouped_process_exporter/pull/17) ([perlun](https://github.com/perlun))

## [v0.5.1](https://github.com/k1LoW/grouped_process_exporter/compare/v0.5.0...v0.5.1) (2019-09-24)

* Fix: Set --group.exclude to grouper [#16](https://github.com/k1LoW/grouped_process_exporter/pull/16) ([k1LoW](https://github.com/k1LoW))

## [v0.5.0](https://github.com/k1LoW/grouped_process_exporter/compare/v0.4.0...v0.5.0) (2019-09-20)

* Add semaphore for restrict number of open files [#15](https://github.com/k1LoW/grouped_process_exporter/pull/15) ([k1LoW](https://github.com/k1LoW))
* Add http server timeout [#14](https://github.com/k1LoW/grouped_process_exporter/pull/14) ([k1LoW](https://github.com/k1LoW))

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
