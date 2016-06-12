# goweb

## Instruction

- Prepare env

```sh
export POSTGRES_USER=postgres
export POSTGRES_PASSWORD=mypostgres
export POSTGRES_DB=users
```

- Start a `postgres` container

```console
docker run -d --name=pq \
    -p 5432:5432 \
    -e POSTGRES_USER=$POSTGRES_USER \
    -e POSTGRES_PASSWORD=$POSTGRES_PASSWORD \
    -e POSTGRES_DB=$POSTGRES_DB \
    postgres
```

- Run the app

```console
go run main.go
```
