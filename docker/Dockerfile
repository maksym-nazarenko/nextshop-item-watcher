FROM golang:1.16-alpine as builder

RUN apk add --no-cache ca-certificates
COPY . /go/src/app/
WORKDIR /go/src/app/
RUN CGO_ENABLED=0 go build -v -o app cmd/main.go


FROM scratch
COPY --from=builder /etc/ssl/ /etc/ca-certificates /etc/ssl/
COPY --from=builder /go/src/app/app /app
ENTRYPOINT [ "/app" ]

USER 10000:10000
