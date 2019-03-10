# slackbot-go
Writing a Slackbot to learn about golang.

Currently, this project contains a HTTP server with 2 routes:

- `/echo`, a `POST` endpoint that returns request body's JSON
- `/`, a `GET` endpoint that displays a message.

## Requirements
- Docker 18.09+
- Docker-Compose 1.24+


# Getting started

To start the server, run
```bash
$ docker-compose up slackbot
```
