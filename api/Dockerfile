# Stage 1: Build Go binary
FROM golang:1.24-alpine AS build
WORKDIR /app

# Install git for go modules fetching
RUN apk add --no-cache git

# Copy go.mod and go.sum first for dependency download
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the code (including local packages)
COPY . .

# Build the application
RUN go build -o /app/main ./cmd/api

# Stage 2: Get wkhtmltopdf binary from a working prebuilt Alpine image
FROM surnet/alpine-wkhtmltopdf:3.19.0-0.12.6-full AS wkhtml

# Stage 3: Final runtime image
FROM alpine:3.21.3 AS prod
WORKDIR /app

# Install CA certificates
RUN apk add --no-cache ca-certificates \
    && update-ca-certificates

# Copy wkhtmltopdf binary & its libraries
COPY --from=wkhtml /bin/wkhtmltopdf /usr/local/bin/wkhtmltopdf
COPY --from=wkhtml /lib /lib
COPY --from=wkhtml /usr/lib /usr/lib

# Create a non-root user and group
RUN addgroup -S openb2b && adduser -S golang -G openb2b

# Copy the Go app binary
COPY --from=build /app/main /app/main

# Change ownership of the application files
RUN chown -R golang:openb2b /app

# Switch to non-root user
USER golang
EXPOSE 3000
CMD ["./main"]
