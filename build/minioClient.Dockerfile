FROM golang:1.21.0-alpine AS builder

# Set the working directory
WORKDIR /ouzi/

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -o ./.bin ./cmd/minio/main.go

FROM scratch AS runner

# Set the working directory
WORKDIR /ouzi/

# Copy the built binary from the builder stage
COPY --from=builder  /ouzi/.bin /ouzi/.bin

# Copy timezone data
COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /
ENV TZ="Europe/Moscow"
ENV ZONEINFO=/zoneinfo.zip

# Expose ports
EXPOSE 8081

# Set the entrypoint
ENTRYPOINT ["./.bin"]