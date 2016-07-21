CREATE KEYSPACE "neurose_order" WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };

USE "neurose_order";

CREATE TABLE "order"(id VARCHAR primary key, reference VARCHAR,price INT,payed INT);

CREATE TABLE "order_item"(id VARCHAR primary key, price INT, quantity INT, deleted BOOLEAN, order_id VARCHAR);

CREATE TABLE "transaction"(id VARCHAR primary key, amount INT, type INT, order_id VARCHAR);