cassandra:
	docker run --name neurose-order -d cassandra:3.0

cassandra-connect:
	docker exec -it neurose-order bash

cassandra-sql:	
	docker run -it --link neurose-order:cassandra --rm cassandra sh -c 'exec cqlsh "$CASSANDRA_PORT_9042_TCP_ADDR"'