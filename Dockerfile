# builder base
FROM golang:1.20-alpine as builder-base
RUN apk update
# install dependency here
RUN apk add build-base ca-certificates && update-ca-certificates tzdata

# create ft user
ENV USER=ft
ENV UID=1000

# See https://stackoverflow.com/a/55757473/12429735RUN 
RUN adduser \    
    --disabled-password \    
    --gecos "" \    
    --home "/nonexistent" \    
    --shell "/sbin/nologin" \    
    --no-create-home \    
    --uid "${UID}" \    
    "${USER}"

# build
FROM builder-base AS builder
RUN mkdir /build
COPY . /build/
WORKDIR /build
RUN go mod download
RUN go mod verify
RUN make test
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s -X main.version=$(cat VERSION) -X 'main.build=$(date)'" -o /build/dist/telegramchatbot

# release base
FROM alpine:3 AS release-base
RUN apk update && apk upgrade
RUN apk add bash
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /
ENV ZONEINFO=/zoneinfo.zip

FROM release-base AS release
COPY --from=builder /build/dist/telegramchatbot /usr/local/bin/
RUN chmod a+rx /usr/local/bin/telegramchatbot
WORKDIR /tmp/
USER ft:ft
ENTRYPOINT [ "telegramchatbot" ]