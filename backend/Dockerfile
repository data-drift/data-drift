FROM golang:1.20
ENV HOME /root
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build


VOLUME $HOME/.datadrift/default
ENV PORT=9740
EXPOSE $PORT

# # Run the application when the container starts
CMD ["./data-drift"]