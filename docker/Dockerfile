FROM alpine:latest
MAINTAINER Mike Petersen <mike@odania-it.com>

ENV TERRALESS_VERSION 0.1.0
ENV TERRALESS_DOWNLOAD_URL https://github.com/Odania-IT/terraless/releases/download/v${TERRALESS_VERSION}/terraless_alpine

# Prepare system
RUN apk add --update terraform bash && rm -rf /var/cache/apk/*
RUN addgroup -g 1000 terraless
RUN adduser -h /srv -G terraless -D -s /bin/bash -u 1000 terraless

# Install terraless & swamp
ADD $TERRALESS_DOWNLOAD_URL /bin/terraless
RUN chmod +x /bin/terraless

WORKDIR /srv
USER terraless

VOLUME /srv

ENTRYPOINT ["/terraless"]
CMD ["-?"]
