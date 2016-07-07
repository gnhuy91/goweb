# goweb [![Build Status](https://travis-ci.org/gnhuy91/goweb.svg?branch=travis-test)](https://travis-ci.org/gnhuy91/goweb)

## Instruction

### Local development

- Install [`glide`](https://github.com/Masterminds/glide) - Go package management tool

- Install dependencies

```console
glide install -s
```

- Prepare env

```sh
export PORT=8080
export POSTGRES_HOST=127.0.0.1:5432
export POSTGRES_USER=postgres
export POSTGRES_PASSWORD=mypostgres
export POSTGRES_DB=users
export UAA_URI=<your_uaa_instance_uri>
```

- Start a `postgres` container

```console
docker run -d --name=pg \
    -p 5432:5432 \
    -e POSTGRES_USER=$POSTGRES_USER \
    -e POSTGRES_PASSWORD=$POSTGRES_PASSWORD \
    -e POSTGRES_DB=$POSTGRES_DB \
    postgres
```

- Building (for *alpine*)

```console
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/goweb
```

- Run the app

```console
bin/goweb
```

### Deploy to Cloud Foundry

- Build the app

```console
go build -o bin/goweb
```

- Push to Cloud Foundry

```console
cf push -f manifest.yml
```
