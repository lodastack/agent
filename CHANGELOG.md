## v0.3.1 [2018-04-18]

### Release Notes

- [#99](https://github.com/lodastack/agent/pull/99): only listen on `localhost` for trace
- [#101](https://github.com/lodastack/agent/pull/101): support IPv6 address report
- [#105](https://github.com/lodastack/agent/pull/105): add agent traffic API

### Bugfixes

- [#100](https://github.com/lodastack/agent/pull/100): fix disk remove bug
- [#104](https://github.com/lodastack/agent/pull/104): close FS monitor file ASAP


## v0.3.0 [2018-02-28]

### Breaking Changes

Add trace config in config file:

```

[trace]
	collector = [ "0.0.0.0:1233" ]

```

if you don't config this , agent will work fun.

### Release Notes

- [#94](https://github.com/lodastack/agent/pull/94): use dep for dependency management
- [#95](https://github.com/lodastack/agent/pull/95): jeager trace support
- [#95](https://github.com/lodastack/agent/pull/95): remove debug subcommand
- [#97](https://github.com/lodastack/agent/pull/97): fix hostname change bug

### Bugfixes

- [#92](https://github.com/lodastack/agent/pull/92): fix disk io await negative value issue

## v0.2.8 [2018-01-04]

### Release Notes

- [#87](https://github.com/lodastack/agent/pull/87): support pcap IPv4 metric collect

### Bugfixes

- [#88](https://github.com/lodastack/agent/pull/88): fix PID file bug
- [#89](https://github.com/lodastack/agent/pull/89): fix output panic
- [#90](https://github.com/lodastack/agent/pull/90): fix pcap data race

## v0.2.7 [2017-12-12]

### Release Notes

- [#85](https://github.com/lodastack/agent/pull/85): support windows serial number collect

### Bugfixes

- [#84](https://github.com/lodastack/agent/pull/84): fix disk io util negative value issue

## v0.2.6 [2017-11-24]

### Release Notes

- [#81](https://github.com/lodastack/agent/pull/81): support report machine serial number

### Bugfixes


## v0.2.5 [2017-11-03]

### Release Notes

- [#79](https://github.com/lodastack/agent/pull/79): update nux lib to display disk block metric

### Bugfixes

- [#81](https://github.com/lodastack/agent/pull/81): ignore docker mount points

## v0.2.4 [2017-09-11]

### Release Notes

- [#74](https://github.com/lodastack/agent/pull/74): support MemAvailabl metric
- [#75](https://github.com/lodastack/agent/pull/75): filter report IP list

### Bugfixes


## v0.2.3 [2017-08-03]

### Release Notes

- [#65](https://github.com/lodastack/agent/pull/65): support collect disk usage blocks metrices

### Bugfixes

- [#67](https://github.com/lodastack/agent/pull/67): add post mq timeout
- [#68](https://github.com/lodastack/agent/pull/68): check plugin report metric name
- [#69](https://github.com/lodastack/agent/pull/69): remove `&` in plugin security check
- [#70](https://github.com/lodastack/agent/pull/70): improve plugin para check

## v0.2.2 [2017-06-09]

### Release Notes

### Bugfixes

- [#61](https://github.com/lodastack/agent/pull/61): provide random seed
- [#62](https://github.com/lodastack/agent/pull/62): fix interface negative value
- [#63](https://github.com/lodastack/agent/pull/63): fix update ip cahce file

## v0.2.1 [2017-05-19]

### Release Notes

### Bugfixes

- [#59](https://github.com/lodastack/agent/pull/59): fix post metric bug
- [#58](https://github.com/lodastack/agent/pull/58): filter windows DVD driver FS

## v0.2.0 [2017-05-18]

### Release Notes

- [#56](https://github.com/lodastack/agent/pull/56): Support windows

### Bugfixes

- [#53](https://github.com/lodastack/agent/pull/53): Don't block if agent queue is full
- [#54](https://github.com/lodastack/agent/pull/54): check agent post metric name

## v0.1.2 [2017-05-02]

### Release Notes

- [#49](https://github.com/lodastack/agent/pull/49): support clean data dir in stop command
- [#51](https://github.com/lodastack/agent/pull/51): Add fs.space.total to sysinfo

### Bugfixes

- [#48](https://github.com/lodastack/agent/pull/48): Update process proc CPU
- [#50](https://github.com/lodastack/agent/pull/50): remove math.Ceil from statics

## v0.1.1 [2017-03-31]

### Release Notes

- [#47](https://github.com/lodastack/agent/pull/47): Support HTTPS scheme

### Bugfixes

- [#44](https://github.com/lodastack/agent/pull/44): Random send report/ns data to server
- [#46](https://github.com/lodastack/agent/pull/46): filter link local multicast addresses

## v0.1.0 [2017-03-10]