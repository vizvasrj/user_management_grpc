##### search could work like this for checking married is false or true and city is Paris

```bash
grpcurl -plaintext -d '{"married": {"is_married": false}, "city": "Paris"}' localhost:50051 user.UserService/SearchUsers
```

# Build and Run the server

```bash
docker compose up --build -d
```

### Test the server

```bash
docker run --name some-postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d grpctemp
```

```bash
go test -v ./...
```
