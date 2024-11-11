# Stage 1: Build the app
FROM golang:latest as builder

RUN mkdir /build 

ADD . /build/

WORKDIR /build 

# Build the Go app with static linking
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Stage 2: Deploy the app
FROM alpine

RUN adduser -S -D -H -h /app appuser

# Copy the app files and binary from the builder stage
COPY . /app

# Copy the statically compiled binary from the builder stage
COPY --from=builder /build/main /app/main

# Set working directory to /app
WORKDIR /app

# Ensure binary has execute permission as root
USER root
RUN chmod +x /app/main

# Switch back to non-root user for running the app
USER appuser

# Expose the port
EXPOSE 3000

# Run the binary
CMD ["./main"]