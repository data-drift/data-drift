# Requirements

In order to run Data Drift, you will need:

- a Github application:
  - It'll need to suubscribe to the events "push".
  - It'll need access to the **content**, read-only
  - When it is created download a secret key \*.private-key.pem
  - Store the github app id as well
  - When your app is running you'll need to update the webhook url to your url + `/webhook/github`

# Hosting

## Heroku

- Create an app and add the buildpack https://github.com/timanovsky/subdir-heroku-buildpack
- Add a config vars `PROJECT_PATH` equals to `backend`
- Add a config vars `GITHUB_APP_PRIVATE_KEY` with your certificate `-----BEGIN RSA PRIVATE KEY-----` (the content of the .private-key.pem)
- Push this repository to heroku

Go to the #Verify section.

## Docker

- Pull the docker image

```
docker pull quay.io/datadrift/data-drift:0.0.1
```

- Run with the GITHUB_APP_PRIVATE_KEY_PATH

```
docker run -v /path/to/local/private-key.pem:/app/private-key.pem -e GITHUB_APP_ID=325270 -e GITHUB_APP_PRIVATE_KEY_PATH=private-key.pem -e GITHUB_APP_ID=github_app_id quay.io/datadrift/data-drift:0.0.1
```

## Kubernetes

Go see [deployment documentation](../self-hosting/k8s/README.md)

# Verify

Go to your URL you should see {"status":"OK"}.
