go-logdna
----

[![GoDoc][1]][2] [![License: MIT][3]][4] [![Release][5]][6] [![Build Status][7]][8] [![Go Report Card][13]][14] [![Code Climate][19]][20] [![BCH compliance][21]][22]

[1]: https://godoc.org/github.com/evalphobia/go-logdna?status.svg
[2]: https://godoc.org/github.com/evalphobia/go-logdna
[3]: https://img.shields.io/badge/License-MIT-blue.svg
[4]: LICENSE.md
[5]: https://img.shields.io/github/release/evalphobia/go-logdna.svg
[6]: https://github.com/evalphobia/go-logdna/releases/latest
[7]: https://travis-ci.org/evalphobia/go-logdna.svg?branch=master
[8]: https://travis-ci.org/evalphobia/go-logdna
[9]: https://coveralls.io/repos/evalphobia/go-logdna/badge.svg?branch=master&service=github
[10]: https://coveralls.io/github/evalphobia/go-logdna?branch=master
[11]: https://codecov.io/github/evalphobia/go-logdna/coverage.svg?branch=master
[12]: https://codecov.io/github/evalphobia/go-logdna?branch=master
[13]: https://goreportcard.com/badge/github.com/evalphobia/go-logdna
[14]: https://goreportcard.com/report/github.com/evalphobia/go-logdna
[15]: https://img.shields.io/github/downloads/evalphobia/go-logdna/total.svg?maxAge=1800
[16]: https://github.com/evalphobia/go-logdna/releases
[17]: https://img.shields.io/github/stars/evalphobia/go-logdna.svg
[18]: https://github.com/evalphobia/go-logdna/stargazers
[19]: https://codeclimate.com/github/evalphobia/go-logdna/badges/gpa.svg
[20]: https://codeclimate.com/github/evalphobia/go-logdna
[21]: https://bettercodehub.com/edge/badge/evalphobia/go-logdna?branch=master
[22]: https://bettercodehub.com/


Unofficial golang library for LogDNA.


# Quick Usage

```go
import (
	"github.com/evalphobia/go-logdna/logdna"
)

func someFunction() {
	conf := logdna.Config{
		APIKey:       "",
		App:          "myapp",
		Env:          "production",
		MinimumLevel: logdna.LogLevelInfo,
		Sync:         false,
		Debug:        true,
	}

	cli, err := logdna.New(conf)
	if err != nil {
		panic(err)
	}

	cli.Debug("logging...")
	cli.Trace("logging...")
	cli.Info("logging...")
	cli.Warn("logging...")
	cli.Err("logging...")
	cli.Fatal("logging...")
}

```


# Environment variables

| Name | Description |
|:--|:--|
| `LOGDNA_API_KEY` | API Key of LogDNA. |
