FROM alpine:edge as builder
LABEL stage=go-builder
RUN apk add --no-cache bash git go gcc musl-dev; \
    sh build.sh docker

FROM alpine:edge
LABEL MAINTAINER="i@nn.ci"
VOLUME ./alist/data/
WORKDIR ./alist/
COPY --from=builder ./bin/alist ./
EXPOSE $PORT
CMD [ "./alist" ]

