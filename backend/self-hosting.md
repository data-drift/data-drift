# Requirements

In order to run Data Drift, you will need:

- a Github application:
  - It'll need to suubscribe to the events "push".
  - It'll need access to the **content**, read-only
  - When it is created keep the secretkey
  - When your app is running you'll need to update the webhook url to your url + `/webhook/github`

# Hosting

## Heroku

- Create an app and add the buildpack https://github.com/timanovsky/subdir-heroku-buildpack
- Add a config vars `PROJECT_PATH` equals to `backend`
- Add a config vars `GITHUB_APP_PRIVATE_KEY` with your certificate `-----BEGIN RSA PRIVATE KEY-----`
- Push this repository to heroku

Go to the #Verify section.

## Docker

@TODO

# Verify

Go to your URL you should see {"status":"OK"}.
