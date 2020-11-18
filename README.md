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

At this time, just a simple notification method was implemented (stdout).
After start the webhook, is possible to make a call and the data will be printed on the stdout.

```bash
curl -i -X POST localhost:3000/api/v0/notifications -d '{"text":"aaa456"}'
```
