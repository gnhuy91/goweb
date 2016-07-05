# goweb

## Instruction
- Install [`glide`](https://github.com/Masterminds/glide) - Go package management tool

- Install dependencies
```console
glide install -s
```

- Prepare env

```sh
export POSTGRES_HOST=127.0.0.1:5432
export POSTGRES_USER=postgres
export POSTGRES_PASSWORD=mypostgres
export POSTGRES_DB=users
export PORT=8080
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

- Run the app

```console
go run *.go
```

- Building (for *alpine*)

```console
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build
```
