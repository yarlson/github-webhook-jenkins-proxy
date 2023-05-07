FROM golang:1.20.1-alpine3.17 as builder

RUN apk --update add ca-certificates && mkdir -p /build
WORKDIR /build
COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags "-w" -trimpath -o /usr/bin/proxy

FROM scratch
ENV PATH=/bin
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/bin/proxy /

ENTRYPOINT ["/proxy"]
