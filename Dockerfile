# Start from base image
FROM golang:1.18-alpine as build

# Set the current working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy source from current directory to working directory
COPY . .

# Build the application
# Produce binary named main
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -a -o main cmd/main.go

########################

# Start from a new lightweight image
FROM alpine:3.18.3

# Install ffmpeg and imagemagick
RUN apk --no-cache add ffmpeg imagemagick

# Set the working directory inside the container
WORKDIR /app

# Copy the binary from the build image to the working directory
COPY --from=build /app/main .

# Run the binary
ENTRYPOINT ["./main"]

# Set default arguments
CMD ["-i", "input.mp4"]
