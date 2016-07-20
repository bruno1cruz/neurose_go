CREATE KEYSPACE "neurose_order" WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };

USE "neurose_order";

CREATE TABLE "order"(id VARCHAR primary key, reference VARCHAR,price INT);


CREATE TABLE "order_item"(id VARCHAR primary key, price INT, quantity INT, order_id VARCHAR);