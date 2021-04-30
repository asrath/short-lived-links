FROM golang:1.15 AS builder

COPY . /go/src/github.com/asrath/short-lived-links

WORKDIR /go/src/github.com/asrath/short-lived-links

RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd

FROM alpine:3.7

COPY --from=builder /go/src/github.com/asrath/short-lived-links/app /opt/sll/app
ADD configs/app.yaml /opt/sll/app.yaml
ADD web /opt/sll/web

RUN mkdir -p /var/sll/pastes
VOLUME ["/var/sll/pastes"]

WORKDIR /opt/sll

EXPOSE 8080
CMD ["./app"]