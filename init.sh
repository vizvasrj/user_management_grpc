#!/bin/sh
set -e

# Wait for PostgreSQL to be ready
until PGPASSWORD=${POSTGRES_PASSWORD} psql -h ${POSTGRES_HOST} -p ${POSTGRES_PORT} -U ${POSTGRES_USER} -d ${POSTGRES_DB} -c "select 1"; do
  >&2 echo "Postgres is unavailable - sleeping"
  sleep 1
done

# Run the migrations
migrate -path ./migrations -database postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@db:5432/$POSTGRES_DB?sslmode=disable up
PGPASSWORD=${POSTGRES_PASSWORD} psql -h ${POSTGRES_HOST} -p ${POSTGRES_PORT} -U ${POSTGRES_USER} -d ${POSTGRES_DB} -f ./seed/seed.sql

# Run app
./main