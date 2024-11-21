# Base Go image
FROM golang:1.23-alpine AS builder

RUN mkdir /app

COPY . /app

WORKDIR /app

# Build the binary without any C libraries
RUN CGO_ENABLED=0 go build -o broker-app ./cmd/api

# Make sure the binary is executable, just in case
RUN chmod +x /app/broker-app

# Build a small image starting from scratch. 
# We just want to have the binary in the final image
FROM alpine:latest

# Create app directory (different from above, this is for the final image)
RUN mkdir /app

# Copy the binary from the builder image into the final image
COPY --from=builder /app/broker-app /app

# Run the app
CMD [ "/app/broker-app" ]
