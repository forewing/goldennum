# Go! Goldennum!

[![CI](https://github.com/forewing/goldennum/workflows/CI/badge.svg)](https://github.com/forewing/goldennum/actions?query=workflow%3ACI)
[![Docker](https://github.com/forewing/goldennum/workflows/Docker/badge.svg)](https://hub.docker.com/r/forewing/goldennum)

## Ideas

[MSRA News: Golden Number Game](https://www.msra.cn/zh-cn/news/features/golden-number-game)

## Rules

There are N players, begin a competition round of Goldennum. Every round, every player submits 2 float numbers in the open interval (0, 100) to the server. At the end of the round, the server will calculate the average of the numbers submitted during this round as AVG. Then the Goldennum of the round is calculated as AVG * 0.618. For every submitted number, if it is the closest number to the Goldennum, the submitter gets N-2 points, if it is the furthest number from the Goldennum, the submitter gets -2 points, otherwise, the submitter gets no points. After the rounds ends, player with the most points wins the game.

## Preview

![preview](https://github.com/forewing/images/raw/master/goldennum-desktop.png)

## Usage

For detail, you may refer to the wiki pages.

- [Deployment Guide](https://github.com/forewing/goldennum/wiki/Deployment)

- [Configuration Guide](https://github.com/forewing/goldennum/wiki/Configuration)

### Build

1. Install [go-bindata](https://github.com/go-bindata/go-bindata)

```
go get -u github.com/go-bindata/go-bindata/go-bindata
```

2. Build resource files

```
go-bindata -fs -prefix "statics/" statics/ templates/
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

Players may develop their own AI/Bot to take part in the game, using the API port.

[API Specification](https://github.com/forewing/goldennum/wiki/API-Spec)

## Development Guide

PR welcome!

[Development Guide](https://github.com/forewing/goldennum/wiki/Development-Guide)

## Thanks

- [Nanjing University Microsoft Student Club](https://github.com/njumsc) for sending me to MSC Summer Camp, 2018, where I met the game.

- [Zhanglv0413](https://github.com/Zhanglv0413) for designing the lovely [logo](./statics/favicon.ico).
