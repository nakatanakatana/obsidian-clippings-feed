FROM --platform=$BUILDPLATFORM golang:1.25 AS builder
ARG TARGETARCH

WORKDIR /app/source
COPY ./ /app/source

ARG CGO_ENABLED=0

RUN mkdir /app/output
RUN GOARCH=${TARGETARCH} go build -o /app/output ./cmd/...

FROM gcr.io/distroless/static-debian12

COPY --from=builder /app/output /app

CMD ["/app/feed"]
