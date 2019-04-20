FROM alpine:edge
MAINTAINER Mike Petersen <mike@odania-it.com>

ENV AWS_PROVIDER=2.7.0
ENV GOOGLE_PROVIDER=2.5.0

# Prepare system
RUN apk add --update --no-cache terraform bash ca-certificates curl && rm -rf /var/cache/apk/*
RUN addgroup -g 1000 terraless
RUN adduser -h /srv -G terraless -D -s /bin/bash -u 1000 terraless

# Install terraless
COPY terraless /bin/terraless

# Preload terraform provider
WORKDIR /bin
RUN curl -LO https://releases.hashicorp.com/terraform-provider-aws/2.7.0/terraform-provider-aws_${AWS_PROVIDER}_linux_amd64.zip \
    && unzip terraform-provider-aws_${AWS_PROVIDER}_linux_amd64.zip \
    && rm terraform-provider-aws_${AWS_PROVIDER}_linux_amd64.zip
RUN curl -LO https://releases.hashicorp.com/terraform-provider-google/2.5.0/terraform-provider-google_${GOOGLE_PROVIDER}_linux_amd64.zip \
    && unzip terraform-provider-google_${GOOGLE_PROVIDER}_linux_amd64.zip \
    && rm terraform-provider-google_${GOOGLE_PROVIDER}_linux_amd64.zip

WORKDIR /srv
USER terraless

VOLUME /srv

ENTRYPOINT ["/bin/terraless"]
CMD ["-?"]
