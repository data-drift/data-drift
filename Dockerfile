FROM node:18 as frontend
WORKDIR /app
COPY frontend/package*.json ./
RUN npm install
COPY frontend/ ./
RUN npm run build:docker


FROM golang:1.20 as backend
ENV HOME /root
WORKDIR /app

COPY backend/go.mod .
COPY backend/go.sum .
RUN go mod download

COPY ./backend .
COPY --from=frontend /app/dist ./dist-app
RUN go build


VOLUME $HOME/.datadrift/default
ENV PORT=9740
EXPOSE $PORT

CMD ["./data-drift"]