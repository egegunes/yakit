db:
	sudo docker run -d -p 5432:5432 --rm --name postgres postgres

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
	go build -o yakitserver

run:
	LISTENADDR=":8080" DBNAME=postgres DBUSER=postgres DBPASS=postgres ./yakitserver

package:
	sudo docker build -t gcr.io/kubernetes-222419/yakit .

push:
	sudo -E docker push eu.gcr.io/kubernetes-222419/yakit
