# name-dyndns [![Build Status](https://travis-ci.org/mfycheng/name-dyndns.svg?branch=master)](https://travis-ci.org/mfycheng/name-dyndns) [![GoDoc](https://godoc.org/github.com/mfycheng/name-dyndns?status.svg)](https://godoc.org/github.com/mfycheng/name-dyndns)

Client that automatically updates name.com DNS records.

## Getting name-dyndns

Since name-dyndns has no external dependencies, you can get it simply by:

```go
go get github.com/mfycheng/name-dyndns.git
```

## Requirements

In order to use name-dyndns, you must have an API key from name.com, which
can be requested from https://www.name.com/reseller/apply.

Once you have your API key, all you must do is setup `config.json`. An example
`config.json` file can be found in `api/config_test.json`.

## Command Line Arguments

By default, running name-dyndns will run a one-time update, using `./config.json`
as a configuration file, and stdout as a log output. However, these can be configured. For example:

```
./name-dyndns -daemon=true -dev=true -log="/var/log/name-dyndns/out.log" -config="~/name_config.json"
```

This will run name-dyndns in daemon mode for dev configurations, outputting to `/var/log/name-dyndns/out.log`, using the configuration file `~/name_config.json`

A detailed usage can be found by running:

```
./name-dyndns --help
```

## Error Handling

Currently, there is limited testing, primarily on none-api dependant utilities.
While error handling _should_ be done gracefully, not every edge case has been tested.

Ideally, when running in daemon mode, name-dyndns tries to treat any errors
arising from network as transient failures, and tries again next iteration. The idea behind this is that a single network failure shouldn't
kill the daemon, which could then potentially result in having the DNS records out
of sync, which would defeat the whole point of name-dyndns.
