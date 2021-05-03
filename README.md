# Webhook Consumer

The purpose of the Webhook Consumer is to provide a easy way to consume the notifications from the Stone services.

## Possible actions when receives a notification

When the Webhook Consumer receives a notification from Stone services, it can take some of the following actions:
* Sends the notification to stdout (just for development)
* Sends the notification to another API
* Stores the notification on a Redis
* Sends the notification to a Kafka topic

## Current state

At this time, this project does not have a stable release.

## Development

To init the development environment and runs the project, follow these steps:

### Clone this repo

```bash
$ git clone git@github.com:stone-co/webhook-consumer.git
```

### Download and install the dependency tools

```bash
$ make setup
```

### Run the tests (optional)

```bash
$ make test
```

### Compile the project

```bash
$ make compile
```

### Run the project

```bash
$ ./build/webhook-consumer
```

## Usage

At this time, just a simple notifier was implemented (stdout).
After start the webhook, is possible to make a call and the data will be printed on the stdout.

### Usage with Docker

First build the Docker Image, or get at Docker Hub.
```bash
$ make build
```

Now create a container with volume to your certificate file, and a environment
variable `PRIVATE_KEY_PATH` to your _.pem_ file.
```bash
$ docker run -v $(pwd)/tests:/usr/share/certificates -e PRIVATE_KEY_PATH="/usr/share/certificates/partner/fakekey.pem" -d stone-co/webhook-consumer:dev
```

Environment variables, and default values:

- PRIVATE_KEY_PATH="tests/partner/fakekey.pem"
- PUBLIC_KEY_PATH="url://https://sandbox-api.openbank.stone.com.br/api/v1/discovery/keys"
- NOTIFIER_LIST=stdout
- API_PORT="3000"
- API_SHUTDOWN_TIMEOUT="5s"

you can pass environment variable with -e flat to docker container run.

```
-e API_PORT="3000"
```
