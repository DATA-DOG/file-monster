# file-monster
Slack bot for cleaning up old files. This is great for free slack teams, where storage is an issue. 

By setting up this bot, users will be able to remove their shared files without any manual clicking.

## Server requirements

MySQL database (for storing access tokens)

## Setup

1. Create slack app
2. Create slack incoming webhook
3. Copy `config.yml.dist` to `config.yml` and fill in required fields
4. Compile project: `go build`
5. Run app: `./file-monster config.yml`
6. Create slash command and point it to server, with port defined in config file

## Usage

Type your command in slack, e.g. `/flushme` to clean old (1 month or older) files

*note: first time each user will have to follow auth instuctions*

or

add extra parameter (`/flushme all`) to clean all public files, which are 1 month or older

