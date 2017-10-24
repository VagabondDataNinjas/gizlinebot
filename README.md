# GIZ Line Bot

## Setup the line bot on your local

Requirements: 

* Golang
* Golang Dependency tool: https://github.com/golang/dep
* Mysql server (see below instructions for importing the schema - assets/init.db)
* Line account with a messaging API provider: https://developers.line.me/console/register/messaging-api/provider/

```
git clone git@github.com:VagabondDataNinjas/gizlinebot.git
cd gizlinebot
dep install

# config file
cp .gizlinebot.example.toml ~/.gizlinebot.toml
# update ~/.gizlinebot.toml with the values for
# GIZLB_LINE_TOKEN and GIZLB_LINE_SECRET from the
# line developer area Messaging API(https://developers.line.me/)
# set correct the config for the SQL parameters

# import the SQL Schema from assets/init.sql

# install ngrok
# https://ngrok.com/download

# start an ngrok tunnel
# the port you give to ngrok should be the same as
# PORT value (default port 8888)
ngrok http 127.0.0.1:8888

# update the line Webhook URL to the ngrok host + "/linewebhook"
eg: https://d2631531.ngrok.io/linewebhook

# start the bot
go run main.go lineBot
```
