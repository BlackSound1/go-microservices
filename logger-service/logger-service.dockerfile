# Base image
FROM alpine:latest

# Create app directory (different from above, this is for the final image)
RUN mkdir /app

# Copy the binary from the builder image into the final image
COPY logger-app /app

# Run the app
CMD [ "/app/logger-app" ]
