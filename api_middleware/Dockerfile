FROM semior/baseimage:latest
LABEL maintainer="Semior <ura2178@gmail.com>"

WORKDIR /srv

ENV GOFLAGS="-mod=vendor"

RUN apk add --no-cache --update postgresql-client && \
    rm -rf /var/cache/apk/*

COPY ./scripts/entrypoint-dev.sh /entrypoint.sh

EXPOSE 2345

CMD ["/entrypoint.sh"]