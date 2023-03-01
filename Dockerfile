# Build stage
FROM golang:latest AS build

# Set the working directory to /go/src/app
WORKDIR /go/src/app

# Copy the source code into the container
COPY . .

# Build the Go app with a static binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o app .

# Final stage
FROM scratch

# Copy the static Go binary into the container
COPY --from=build /go/src/app/app /app

# Expose port 8080
EXPOSE 8080

# Set the entry point of the container to the Go app
ENTRYPOINT ["/app"]