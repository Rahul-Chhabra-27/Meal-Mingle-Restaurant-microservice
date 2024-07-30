# Base image
FROM golang:1.22.2-alpine3.19

# Install bash
RUN apk add --no-cache bash

# Move to working directory /app
WORKDIR /app

# Copy the code into the container
COPY . .

# Download dependencies using go mod
RUN go mod tidy && go mod vendor

# Expose PORT 8091 to the outside world
EXPOSE 8091

# Command to run the application when starting the container
CMD ["go", "run", "."]
