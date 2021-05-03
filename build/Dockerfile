FROM ubuntu:20.04

RUN apt-get update -y \
    && DEBIAN_FRONTEND=noninteractive apt-get install locales tzdata openssl ca-certificates -y \
    && /usr/sbin/update-ca-certificates \
    && locale-gen en_US.UTF-8 pt_BR.UTF-8

ENV	LANG=en_US.UTF-8 LANGUAGE=en_US:en LC_ALL=en_US.UTF-8 DEBIAN_FRONTEND=noninteractive

ADD webhook-consumer /bin/

ENTRYPOINT ["/bin/webhook-consumer"]
