FROM golang:1.24-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

COPY . .

# Build the application
RUN go build -o /app/main ./cmd/api


# Stage 2
FROM alpine:3.21.3 AS prod

WORKDIR /app

# Create a non-root user and group
RUN addgroup -S openb2b && adduser -S golang -G openb2b

COPY --from=build /app/main /app/main

# Change ownership of the application files
RUN chown -R golang:openb2b /app

# Switch to non-root user
USER golang

EXPOSE 3000

CMD ["./main"]