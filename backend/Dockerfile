# Use an official Golang runtime as a parent image
FROM golang:1.20

# Set the working directory to /app
WORKDIR /app

# Copy the go.mod and go.sum files
COPY go.mod .
COPY go.sum .

RUN go install

# Copy the current directory contents into the container at /app
COPY . .

# Build the Go binary
RUN go build

# Expose port 8080 for the application
EXPOSE 8080

# Run the application when the container starts
CMD ["./data-drift"]