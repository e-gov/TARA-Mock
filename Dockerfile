ARG DOCKERHUB_MIRROR
FROM ${DOCKERHUB_MIRROR}library/golang:1.15.1-alpine

ARG ALPINE_MIRROR
RUN if [ -n "${ALPINE_MIRROR}" ]; then \
    echo "${ALPINE_MIRROR}v3.12/main/" > /etc/apk/repositories && \
    echo "${ALPINE_MIRROR}v3.12/community/" >> /etc/apk/repositories; fi

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