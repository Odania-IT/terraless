FROM alpine:edge
MAINTAINER Mike Petersen <mike@odania-it.com>

# Prepare system
RUN apk add --update --no-cache terraform bash ca-certificates && rm -rf /var/cache/apk/*
RUN addgroup -g 1000 terraless
RUN adduser -h /srv -G terraless -D -s /bin/bash -u 1000 terraless

# Install terraless
COPY terraless /bin/terraless

WORKDIR /srv
USER terraless

VOLUME /srv

ENTRYPOINT ["/bin/terraless"]
CMD ["-?"]
