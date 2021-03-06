FROM golang:1.17-alpine

ARG SERVICE_ENV=dev
ENV GOOS=$TARGETOS
ENV GOARCH=$TARGETARCH

RUN apk add git

COPY . ./build
WORKDIR ./build

RUN go mod download && \
    go mod verify

RUN scripts/build.sh

FROM alpine:latest
COPY --from=0 /go/build/dist/otter-gateway /usr/bin/otter-gateway
COPY --from=0 --chown=1001:1001 /go/build/web /var/lib/otter/web

RUN adduser --system --disabled-password --shell /bin/sh --home /var/lib/otter --uid 1001 otter && \
    addgroup --system --gid 1001 otter && \
    mkdir -p /var/lib/otter; chown otter:otter -R /var/lib/otter

USER otter
WORKDIR /var/lib/otter
EXPOSE 50051
CMD ["/usr/bin/otter-gateway"]
