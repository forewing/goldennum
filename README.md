# Go! Goldennum!

## Usage

### Config

1. Default config file is `conf/config.yml`, you may rename `conf/config.example.yml` and edit it as you prefer.

2. You may override the default config file path by setting environment variable `CONFIG_PATH`

3. You may also use environment variable as config by setting `USE_ENV_CONFIG=true`, see `example.env` for detail.

### Build

1. Use `build.sh` to produce binary as `output/goldennum`. The script is as simple as running `go build` and move the binary to `output/`.

### Running

1. `run.sh` will export all environment variable from `.env`, you may rename `example.env` and edit it as you prefer.

## Database

> NOTICE: default settings use in memory sqlite, which means you will lost all of the data after restart the program. You should change sqlite path to file path or use mysql instead for data persistence.

## API

TODO
