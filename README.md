# Go! Goldennum!

[![CI](https://github.com/forewing/goldennum/workflows/CI/badge.svg)](https://github.com/forewing/goldennum/actions?query=workflow%3ACI)
[![Docker](https://github.com/forewing/goldennum/workflows/Docker/badge.svg)](https://hub.docker.com/r/forewing/goldennum)

## Ideas

[MSRA News: Golden Number Game](https://www.msra.cn/zh-cn/news/features/golden-number-game)

## Rules

There are N players, begin a competition round of Goldennum. Every round, every player submits 2 float numbers in the open interval (0, 100) to the server. At the end of the round, the server will calculate the average of the numbers submitted during this round as AVG. Then the Goldennum of the round is calculated as AVG * 0.618. For every submitted number, if it is the closest number to the Goldennum, the submitter gets N-2 points, if it is the furthest number from the Goldennum, the submitter gets -2 points, otherwise, the submitter gets no points. After the rounds ends, player with the most points wins the game.

For players, take a look at our [Participation Guide](https://github.com/forewing/goldennum/wiki/Participation-Guide)!

> If you want to host a game but having trouble deploying on your server, feel free to contact me. I am glad to open a lobby on my server for you.

## Preview

![preview](https://github.com/forewing/images/raw/master/goldennum-desktop.png)

## Usage

For detail, you may refer to the wiki pages.

- [Deployment Guide](https://github.com/forewing/goldennum/wiki/Deployment)

- [Configuration Guide](https://github.com/forewing/goldennum/wiki/Configuration)

## Install

### Docker

You can get the docker image at [forewing/goldennum](https://hub.docker.com/r/forewing/goldennum).

See [docker-compose.yml](docker-compose.yml) for example config.

### Pre-Built Binary

Download it from [GitHub Action CI](https://github.com/forewing/goldennum/actions?query=workflow%3ACI+is%3Asuccess), latest successful build is recommended.

Or from [Latest Release](https://github.com/forewing/goldennum/releases/tag/v0.3.8) page.

### Build From Source

> Require Go 1.16+

Make sure you have `GO111MODULE` set to `on`

```
go env -w GO111MODULE=on
```

> You may need to setup [GOPROXY](https://github.com/goproxyio/goproxy) if you live in countries without international network connections.
>
> To set it up, simply run `go env -w GOPROXY=https://goproxy.io,direct`.

You can install with a single command

```
go install github.com/forewing/goldennum
```

Or get the source code and build manually

```
git clone https://github.com/forewing/goldennum.git
cd goldennum
go build
```

## Run

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

[API Specification](https://github.com/forewing/goldennum/wiki/API-Specification)

## Development Guide

PR welcome!

[Development Guide](https://github.com/forewing/goldennum/wiki/Development-Guide)

## Thanks

- [Nanjing University Microsoft Student Club](https://github.com/njumsc) for supporting this project!

- [Zhanglv0413](https://github.com/Zhanglv0413) for designing the lovely logo!

    [![logo](resources/statics/favicon.ico)](resources/statics/favicon.ico)
