# Requirements

In order to run Data Drift, you will need:

- Create a [new Github application](https://docs.github.com/en/apps/creating-github-apps/registering-a-github-app/registering-a-github-app)
- The GH application parameters:
  - In the homepage URL set: https://datadrift.yourdomain.com/
  - In the webhook URL set: https://datadrift.yourdomain.com/webhook/github
  - In the permissions:
    - It'll need access to the **content**, read-only
    - It'll need to subscribe to the events "push".
- When it is created download a secret key \*.private-key.pem
- Store the github app id as well

When the app is created, click "Public page" and install the app on your datadrift repository.

# Hosting

## Heroku

- Create an app and add the buildpack https://github.com/timanovsky/subdir-heroku-buildpack
- Add a config vars `PROJECT_PATH` equals to `backend`
- Add a config vars `GITHUB_APP_PRIVATE_KEY` with your certificate `-----BEGIN RSA PRIVATE KEY-----` (the content of the .private-key.pem)
- Add a config vars `GITHUB_APP_ID` with you gihub app id
- Push this repository to heroku

Go to the #Verify section.

## Docker

- Pull the docker image

```
docker pull quay.io/datadrift/data-drift:0.0.1
```

- Run with the GITHUB_APP_PRIVATE_KEY_PATH

```
docker run -v /path/to/local/private-key.pem:/app/private-key.pem -e GITHUB_APP_ID=your_app_id -e GITHUB_APP_PRIVATE_KEY_PATH=private-key.pem -e GITHUB_APP_ID=github_app_id quay.io/datadrift/data-drift:0.0.1
```

## Kubernetes

Go see [deployment documentation](../self-hosting/k8s/README.md)

# Verify

Go to your URL you should see {"status":"OK"}.
Go to /ghhealth you should see {"status":"OK"}.

Make sure you have updated the webhook url to be /webhook/github.
Go to the "Advanced" section of your github app, and redeliver the last webhook, it should succeed.
