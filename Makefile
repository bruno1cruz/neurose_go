cassandra:
	docker run --name neurose-order -d cassandra

scylladb:
	docker run -d -p 9042:9042 -i -t --name neurose-scylladb scylladb/scylla

cassandra-connect:
	docker exec -it neurose-order bash

cassandra-sql:	
	docker run -it --link neurose-order:cassandra --rm cassandra sh -c 'exec cqlsh "$CASSANDRA_PORT_9042_TCP_ADDR"'

docker-app:
	docker build -t neurose-go .

docker-run:
	docker run -P -it --rm --name neurose-go neurose-go

docker-build:
	docker run --rm -v "$PWD":/usr/src/myapp -w /usr/src/myapp golang:1.6 go build -v