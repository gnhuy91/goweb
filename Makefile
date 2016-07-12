IMAGE=gnhuy91/goweb
DCR=docker-compose run --rm

POSTGRES_HOST?=127.0.0.1:5432
POSTGRES_USER?=postgres
POSTGRES_PASSWORD?=mypostgres
POSTGRES_DB?=users

.PHONY: clean test build release docker-build docker-push run

all: release

clean:
	rm -f bin/*

test:
	$(DCR) go-test

build:
	$(DCR) go-build

run:
	POSTGRES_USER=$(POSTGRES_USER) \
	POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
	POSTGRES_DB=$(POSTGRES_DB) \
	bin/goweb

release: test build docker-build docker-push

docker-build:
	docker build --rm -t $(IMAGE) .

docker-push:
	docker push $(IMAGE)
