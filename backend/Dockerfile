FROM golang:1.20

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build

# # Expose port 8080 for the application
EXPOSE 8080

# # Run the application when the container starts
CMD ["./data-drift"]