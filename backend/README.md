# Quickstart

```sh
docker volume create driftdb-store
docker run --name driftdb -p 9740:9740 -vdriftdb-store:/root/.datadrift quay.io/datadrift/driftdb
```

# Development

In the `backend` folder build the image,

```sh
docker volume create driftdb-store
docker build -t driftdb-dev .
docker run -p 9740:9740 -v driftdb-store:/root/.datadrift driftdb-dev
```
