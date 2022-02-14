FROM docker.io/golang:latest
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
ENTRYPOINT []
CMD ["/app/do-kyoka"]
