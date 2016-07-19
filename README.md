# goweb [![Build Status](https://travis-ci.org/gnhuy91/goweb.svg?branch=develop)](https://travis-ci.org/gnhuy91/goweb) [![GoDoc](https://godoc.org/github.com/gnhuy91/goweb?status.svg)](http://godoc.org/github.com/gnhuy91/goweb)

## Instruction

### Local development

- Run tests

```console
make test
```

- Run the app (in docker container), if you prefer to run the app without docker, skip to the next step

```console
make run-docker
```

- Build & run the app's binary, this requires a **PostgreSQL** instance

```sh
# Build app binary (this output `goweb` binary in `bin` folder)
make build

# Prepare env
export POSTGRES_HOST=127.0.0.1:5432
export POSTGRES_USER=postgres
export POSTGRES_PASSWORD=mypostgres
export POSTGRES_DB=users
export PORT=8080  # goweb's listen port

# Start a postgres container (or spin up your own instance here)
docker run -d --name=pg \
    -p 5432:5432 \
    -e POSTGRES_USER=$POSTGRES_USER \
    -e POSTGRES_PASSWORD=$POSTGRES_PASSWORD \
    -e POSTGRES_DB=$POSTGRES_DB \
    postgres

# Run the binary
make run
```

### Deploy to Cloud Foundry

```console
make cf
```
