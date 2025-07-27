# Stage 1: Build the Go binary
FROM golang:1.22 as builder

WORKDIR /app

# Copy go mod and download dependencies
COPY go.mod ./
RUN go mod download

# Copy source code
COPY . .

# Build the Go binary
RUN go build -o main main.go

# Stage 2: Create a minimal runtime image
FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 8080

# Run the binary
CMD ["/app/main"]
