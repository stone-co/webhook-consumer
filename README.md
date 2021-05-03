# Webhook Consumer

The purpose of the Webhook Consumer is to provide a easy way to consume the notifications from the Stone services.

<img src="/docs/overview.png" alt="Image of webhook consumer diagram" align="center" />

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

Notifier List:

- stdout
- http proxy
- redis

_You can implements a notifier and submit a Pull Request, or showing interest by
creating an issue with the type of notifier_

### Setup

Define `PORT` environment variable to Open Banking Organization send messages to
Webhook Consumer, and customize shutdown timeout with `API_SHUTDOWN_TIMEOUT`.
The defaults values are _3000_ and _5s_.

The environment variable `PRIVATE_KEY_PATH` contains a path to your key file,
your private key made to Open Banking Partner, and `PUBLIC_KEY_PATH` identify
the location of public key from Open Banking Organization.

The environment variable `NOTIFIER_LIST` must be a string, with notifiers name
separated by `;` character.

```bash
$ NOTIFIER_LIST="stdout;proxy;redis"
```

If you use **http proxy** as a notifer you must set the following environment
variables:

- PROXY_NOTIFIER_URL
- PROXY_NOTIFIER_TIMEOUT _(default = 10s)_

If you use **redis** as a notifer you must set the following environment
variables:

- REDIS_ADDR _required_
- REDIS_PORT _required_
- REDIS_PASSWORD
- REDIS_USE_TLS _default false_
- REDIS_MAX_IDLE _default 100_
- REDIS_MAX_ACTIVE _default 1000_
- REDIS_IDLE_TIMEOUT _default 1m_
- REDIS_CONNECT_TIMEOUT _default 1s_
- REDIS_READ_TIMEOUT _default 300ms_
- REDIS_WRITE_TIMEOUT _default 300ms_

Check configure notifer files to view all environment variables:

- [proxy http](/pkg/gateways/notifiers/proxy/configure.go)
- [redis](/pkg/gateways/notifiers/redis/config.go)


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
