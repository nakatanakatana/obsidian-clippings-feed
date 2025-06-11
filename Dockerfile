FROM golang:1.24 AS builder

WORKDIR /app/source
COPY ./ /app/source

ARG CGO_ENABLED=0
ARG GOOS=linux
ARG GOARCH=amd64

RUN mkdir /app/output
RUN go build -o /app/output ./cmd/...

FROM gcr.io/distroless/static-debian12

COPY --from=builder /app/output /app

CMD ["/app/feed"]
