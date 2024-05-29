# Use an official Golang runtime as a parent image
FROM golang:1.21.0 AS builder

# Set the working directory in the container
WORKDIR /app

# Copy the local package files to the container's workspace
COPY . .

# Build the Go app
RUN go build -o main ./cli/user_service/main.go

# Use postgres as base image for the runtime
FROM postgres:latest

# Set the working directory
WORKDIR /docker-entrypoint-initdb.d

# Copy the migrations folder and seed data from the builder stage
COPY --from=builder /app/migrations /docker-entrypoint-initdb.d/
COPY --from=builder /app/seed/seed.sql /docker-entrypoint-initdb.d/
