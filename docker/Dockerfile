# Use an official Golang runtime as a parent image
FROM golang:1.21.0 AS builder

# Set the working directory in the container
WORKDIR /app

# Copy the local package files to the container's workspace
COPY . .

# Install migrate tool
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Build the Go app
RUN go build -o main ./cli/user_service/main.go

# Use a minimal image for the runtime
FROM alpine:latest

# Install necessary dependencies
RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /root/

# Copy the built binary and other necessary files from the builder stage
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate
COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations

# Run the application
CMD ["./main"]
