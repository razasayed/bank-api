CREATE DATABASE bank_db_test;

\connect bank_db_test

\i /docker-entrypoint-initdb.d/schema.sql