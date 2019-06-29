TAG=yakit

build:
	sudo docker build -t $(TAG) .
up:
	sudo docker run --rm --env-file .env --network yakit -p 8080:8080 $(TAG)
rundb:
	sudo docker run -d --rm --env-file .env --name yakitdb --network yakit postgres
psql:
	sudo docker exec -it postgres psql -U postgres postgres
provision:
	scripts/provision.sh
