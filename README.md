# GIZ Line Bot (App Engine)

## Setting up
### Install and initialize Google Cloud SDK
**You can skip this part if you prefer to use [Google Cloud Shell](https://cloud.google.com/shell/docs/quickstart) instead**

1. Follow instruction [here](https://cloud.google.com/appengine/docs/standard/go/download)
2. If you have installed the Google Cloud SDK on your machine, run the command `gcloud auth application-default login`

### Create your new yaml file
```
cd myapp/
cp .local-example.yaml local.yaml
cp .cloud-example.yaml cloud.yaml
```
### Config yaml files
- Specify your `service` name in both yaml files
- Update your `CHANNEL_SECRET` and `CHANNEL_TOKEN` from the line developer area [Messaging API](https://developers.line.me/)

## Run the line bot locally
### Setup mysql server on your local machine
1. Install via [Homebrew](https://gist.github.com/nrollr/3f57fc15ded7dddddcc4e82fe137b58e) or Download and install mysql community server [here](https://dev.mysql.com/downloads/mysql/)
2. Import the SQL Schema from assets/init.sql

### Setup ngrok
Install [ngrok](https://ngrok.com/download)

### Start the bot
1. Config SQL parameters in local.yaml if needed (Default host: 127.0.0.1, port: 3306)
2. Start an ngrok tunnel, run `ngrok http 127.0.0.1:8888`
3. start the bot
    - Make sure you've **installed and started mysql server** on your local machine
    - Run `dev_appserver.py local.yaml --port 8888`

## Deploy the line bot on app engine
1. Config SQL parameters in cloud.yaml if needed
2. Run `gcloud app deploy cloud.yaml --version YOUR-VERSION`
