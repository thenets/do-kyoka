FROM docker.io/golang:latest as builder
WORKDIR $GOPATH/src/github.com/thenets/do-kyoka
RUN set -x \
    && apt-get update \
    && apt-get install -y --no-install-recommends \
        git
COPY . .
ENV GO111MODULE auto
RUN set -x \
    && go get -d -v ./... \
    && go build -ldflags="-extldflags=-static" -o /app/do-kyoka \
    && chmod +x /app/do-kyoka

FROM docker.io/ubuntu:latest
WORKDIR /app
RUN set -x \
    && apt-get update \
    && apt-get upgrade -y --auto-remove \
    && apt-get install -y make \
    && apt-get clean
COPY --from=builder /app/* .
ENTRYPOINT []
CMD ["/app/do-kyoka"]
