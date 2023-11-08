# Development

In the `backend` folder build the image,

```sh
docker volume create datadrift_store # A volume should be created to keep the stored data between container restarts.
docker build -t datadrift .
docker run -p 8080:8080 -v datadrift_store:/root/.datadrift datadrift
```
