# SET UP

## Set Up the DB - Postgres db

        CREATE DATABASE demo_db;

        CREATE USER demo_user WITH SUPERUSER CREATEDB CREATEROLE LOGIN PASSWORD 'demo_pass';

        ALTER ROLE demo_user SET client_encoding TO 'utf8';

        ALTER ROLE demo_user SET default_transaction_isolation TO 'read committed';

        ALTER ROLE demo_user SET timezone TO 'UTC';

        GRANT ALL PRIVILEGES ON DATABASE demo_db TO demo_user;

        CREATE DATABASE demo_test;

        GRANT ALL PRIVILEGES ON DATABASE demo_test TO demo_user;

        \q

## Setting up the project

Clone the repository and if you have go set up on your machine run

        go mod download

        go mod tidy

jwahome@emtechhouse.co.ke
