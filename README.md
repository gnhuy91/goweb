# goweb [![Build Status](https://travis-ci.org/gnhuy91/goweb.svg?branch=travis-test)](https://travis-ci.org/gnhuy91/goweb)

## Instruction

### Local development

- Run tests

```console
make test
```

- Build app binary (this output `goweb` binary in `bin/`)

```console
make build
```

- Run the app

```sh
# Prepare env
export POSTGRES_HOST=127.0.0.1:5432
export POSTGRES_USER=postgres
export POSTGRES_PASSWORD=mypostgres
export POSTGRES_DB=users
export PORT=8080  # goweb's listen port

# Start a postgres container
docker run -d --name=pg \
    -p 5432:5432 \
    -e POSTGRES_USER=$POSTGRES_USER \
    -e POSTGRES_PASSWORD=$POSTGRES_PASSWORD \
    -e POSTGRES_DB=$POSTGRES_DB \
    postgres

# Run the app
make run
```

### Deploy to Cloud Foundry

- Build app binary

```console
make build
```

- Push to Cloud Foundry

```console
cf push -f manifest.yml
```
