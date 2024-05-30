# GRPC User Service

This guide provides instructions for installing dependencies, building and running the server, testing the server, and managing migrations for a GRPC User Service.

## You can use postgres as well as inmemory database just change `.env` with `USE_DATABASE=inmemory` or `USE_DATABASE=postgres`

#### for `inmemory` database you don't need to run migrations and do postgres setup

## Installation

### Install dependencies for interacting with the GRPC server

To install `grpcurl`, use the following command:

```bash
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```

## Build and Run the Server

To build and run the server using Docker, execute the following command:

```bash
docker compose up --build -d
```

### Stop the Server

To stop the server, use the following command:

```bash
docker compose stop
```

### Remove the Server

To remove the server, use the following command:

```bash
docker compose down
```

# Test in Local Development Server

To run the server with a PostgreSQL database, use the following command:

```bash
docker run -it --name some-postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=grpctemp \
  -p 5432:5432 \
  --sig-proxy=true \
  -v postgres_data:/var/lib/postgresql/data postgres:16
```

## Install Migration Dependencies

To install dependencies for database migrations, use the following command:

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

## Install Environment Variable Management Dependencies

To install dependencies for managing environment variables, use the following command:

```bash
go install github.com/joho/godotenv/cmd/godotenv@latest
```

## Use the .env File

To load environment variables from the `.env` file, use the following command:

```bash
source .env
```

## Run Migrations Before Testing

Before testing, run database migrations with the following command, so i use `source .env` here for getting `$POSTGRES_USER` and other variable from .env file:

```bash
migrate -path ./migrations -database postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@localhost/$POSTGRES_DB?sslmode=$POSTGRES_SSLMODE up
```
## Add some data to database

```bash
PGPASSWORD=${POSTGRES_PASSWORD} psql -h ${POSTGRES_HOST} -p ${POSTGRES_PORT} -U ${POSTGRES_USER} -d ${POSTGRES_DB} -f ./seed/seed.sql
```


## Run Tests

To run tests with the environment variables loaded, use the following command:

```bash
godotenv -f .env go test -v ./...
```

## Clean Local Docker Environment

To clean your local Docker environment by removing the container and volume, use the following commands:

```bash
docker container stop some-postgres
docker container rm some-postgres
docker volume rm postgres_data
```

---

### Getting user from user id

You can put id of user in the following command to get the user:

```bash
grpcurl -plaintext -d '{"id": 1}' localhost:50051 user.UserService/GetUserById
# {
#   "id": 1,
#   "fname": "Jane",
#   "city": "Tokyo",
#   "phone": "1764788495",
#   "height": 4.9,
#   "married": true
# }

```

### Getting all users with list of ids

you can put list of ids in the following command to get the users:

```bash
grpcurl -plaintext -d '{"ids": [1, 3]}' localhost:50051 user.UserService/GetUsersByIds
# {
#   "users": [
#     {
#       "id": 1,
#       "fname": "Jane",
#       "city": "Tokyo",
#       "phone": "1764788495",
#       "height": 4.9,
#       "married": true
#     },
#     {
#       "id": 3,
#       "fname": "Jane",
#       "city": "Sydney",
#       "phone": "5660612778",
#       "height": 5.2
#     }
#   ]
# }

```

### Checking if a user is married and living in Paris

You can use the following command to search for users with specific criteria:

```bash
grpcurl -plaintext -d '{"married": {"is_married": false}, "city": "Paris"}' localhost:50051 user.UserService/SearchUsers
# {
#   "users": [
#     {
#       "id": 4,
#       "fname": "Olivia",
#       "city": "Paris",
#       "phone": "7499978875",
#       "height": 5
#     },
#     {
#       "id": 8,
#       "fname": "Susan",
#       "city": "Paris",
#       "phone": "5307397290",
#       "height": 6.1
#     }
#   ]
# }
```

### You can search by height range and city as well:

```bash
grpcurl -plaintext -d '{"city": "Madrid", "height": {"start_value": 3.2, "end_value": "6.2"}}' localhost:50051 user.UserService/SearchUsers
# {
#   "users": [
#     {
#       "id": 9,
#       "fname": "Olivia",
#       "city": "Madrid",
#       "phone": "5199527895",
#       "height": 4.6
#     },
#     {
#       "id": 10,
#       "fname": "Olivia",
#       "city": "Madrid",
#       "phone": "8354340417",
#       "height": 5.5,
#       "married": true
#     }
#   ]
# }
```

### You can search by phone number, hight range, city, fname and married status as well
