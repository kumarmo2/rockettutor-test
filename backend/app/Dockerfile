FROM golang:1.23.4-bookworm as build

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -o api .


FROM debian:bookworm-slim as final

WORKDIR /app

COPY --from=build /app/api .

ENTRYPOINT [ "/app/api" ]
