#!/bin/sh

cd /srv

go version

# run debug
if [ ! -z "${DEBUG}" ]; then
    /go/bin/dlv debug --headless --log --listen=:2345 --api-version=2 /srv/app
	if [ $? -ne 0 ]; then
        exit 1
    fi

fi

go build -mod=vendor -o /go/build/app ./app
/go/build/app