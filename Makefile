.PHONY: doc

TAG=yakit

all: build up

build:
	sudo docker build -t $(TAG) .
up:
	sudo docker run --rm --env-file .env --network yakit -p 8080:8080 $(TAG)
stopdb:
	sudo docker stop yakitdb || true
rundb:
	sudo docker run -d --rm --env-file .env --name yakitdb --network yakit postgres
psql:
	sudo docker exec -it yakitdb psql -U yakit yakit
provision:
	scripts/provision.sh
doc:
	sudo docker run --rm -it -v $(PWD)/doc:/src/doc egegunes/redoc
