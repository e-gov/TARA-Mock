FROM golang:1.15.1-alpine

RUN apk update && apk upgrade && \
    apk add --no-cache bash openssl

COPY ./service /service
COPY ./genkeys /genkeys

WORKDIR /genkeys
RUN chmod +x ./genkeys.sh

ARG genkeys
RUN if [ "$genkeys" = "true" ] ; then ./genkeys.sh; fi

WORKDIR /service
CMD ["go", "run", ".", "-conf", "./config.json"]