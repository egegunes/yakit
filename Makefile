.PHONY: all db initial build run logs dummy

db:
	sudo docker run -d -p 5432:5432 --network yakit --name postgres postgres

initial:
	sudo docker cp initial.sql postgres:/initial.sql
	sudo docker exec postgres psql -U postgres -f /initial.sql postgres

psql:
	sudo docker exec -it postgres psql -U postgres postgres

dummy:
	for file in $(CURDIR)/dummy/*; do \
	    name=$$(basename $$file); \
	    sudo docker cp $$file postgres:/$$name; \
	    sudo docker exec postgres psql -U postgres -f /$$name postgres; \
	done; \

build:
	sudo docker build -t yakit .

run:
	sudo docker build -t yakit .
	sudo docker stop yakit || true
	sudo docker run --rm --network yakit -p 8080:8080 \
			-e LISTENADDR=":8080" \
			-e DBHOST=postgres \
			-e DBNAME=postgres \
			-e DBUSER=postgres \
			-e DBPASS=postgres \
			--name yakit \
			yakit

logs:
	sudo docker logs --follow yakit

gcr-tag:
	sudo docker tag yakit eu.gcr.io/kubernetes-222419/yakit

push:
	sudo -E docker push eu.gcr.io/kubernetes-222419/yakit
