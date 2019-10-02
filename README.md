# paris-pound-check
Check if your vehicle 🏍️ has been impounded by Paris city agents 👮🏻‍♂️

## Installation

Easiest way to use this is to compile it in a docker container :
```
docker build -t paris-pound-check:latest -f ./Dockerfile ./
```

## Running paris-pound-check

Then, running the script itself is a matter of populating environment variables in `.env` and starting docker container with environment file mounted

```
docker run --rm -v $(pwd)/.env:/go/bin/.env paris-pound-check:latest --notifier=slack check
```

The official image built from this repository is also available on Docker hub image registry [https://cloud.docker.com/u/pauulog/repository/docker/pauulog/paris-pound-check](https://cloud.docker.com/u/pauulog/repository/docker/pauulog/paris-pound-check)

### Cron usage

In the following example, the program will check every minute if the vehicle has been impounded.
```
* * * * docker run --rm -v $(pwd)/.env:/go/bin/.env docker.io/pauulog/paris-pound-check:latest --notifier=slack check
```

## Available notifiers

Right now the program can notify you using Slack only. In the future a twilio SMS will maybe be added.
