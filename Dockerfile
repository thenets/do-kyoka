# Build
FROM golang:1.15-alpine
WORKDIR $GOPATH/src/github.com/thenets/do-kyoka
RUN apk add git
COPY . .
RUN go get -d -v ./...
RUN go build -o /tmp/do-kyoka
RUN chmod +x /tmp/do-kyoka

# Server
FROM alpine
ENV FIREWALL_NAME=
RUN adduser -S -D -H -h /app/ kyoka
USER kyoka
WORKDIR /app/
COPY --from=0 /tmp/do-kyoka ./
CMD ["./do-kyoka"]
