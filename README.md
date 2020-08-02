# Go! Goldennum!

[![CI](https://github.com/forewing/goldennum/workflows/CI/badge.svg)](https://github.com/forewing/goldennum/actions?query=workflow%3ACI)
[![Docker](https://github.com/forewing/goldennum/workflows/Docker/badge.svg)](https://hub.docker.com/r/forewing/goldennum)

## Preview

![preview](https://github.com/forewing/images/raw/master/goldennum-desktop.png)

## Usage

### Build

1. Install [packr2](https://github.com/gobuffalo/packr/tree/master/v2)

```
go get github.com/gobuffalo/packr/v2/packr2
```

2. Build resource files

```
packr2
```

3. Build project

```
go build
```

### Run

```
./goldennum -h
Usage of ./goldennum:
  -admin string
        Admin username. (default "admin")
  -base-url string
        Base URL. If you are using reverse proxy to redirect "//PUBLIC_HOST/PREFIX/uri" to "//REAL_HOST/url", it should be set to "/PREFIX"
  -bind string
        Bind address. (default "0.0.0.0:8080")
  -conf string
        Config file path. If set, load config from file instead.
  -db-addr string
        Mysql server address. (default "localhost:3306")
  -db-name string
        Database name. (default "goldennum")
  -db-pass string
        Database password. (default "goldennum")
  -db-path string
        Path to sqlite3 database. (default "./sqlite3.db")
  -db-type string
        [sqlite3, mysql] (default "sqlite3")
  -db-user string
        Database username. (default "goldennum")
  -debug
        Set debug mode.
  -password string
        Admin password. Random if empty.
```

1. By default, server read all configs from command line flags.

2. If flag `-conf={FILE}` is set, server will load configs from `{FILE}` instead of flags. Refer to `conf/config.example.yml` for detail.

3. If environment variable `USE_ENV_CONFIG=true` is set, non-empty environment variable will override corresponding fields. Refer to `example.env` for detail.

4. If `password` not set, a safe random token will be used.

## API

TODO

## Development Guide

- When `conf.Debug` is set to `true`, and the program can access `./templates` & `./statics`, it will live reload those files.
