# syntax=docker/dockerfile:1.4.1
FROM golang:1.19.0-alpine3.15@sha256:1757709b9c057600291a1265d53d39a86a72174fe267b6aa9e2983babf8e9282 AS go
WORKDIR /app
RUN addgroup --gid 99 --system haproxy && \
    adduser --disabled-password --home /var/lib/haproxy --ingroup haproxy --no-create-home --system --uid 99 haproxy && \
    mkdir /var/lib/haproxy && \
    mkdir /data && \
    chown haproxy:haproxy /var/lib/haproxy && \
    chown -R haproxy:haproxy /data && \
    chown -R haproxy:haproxy /app


FROM go AS builder
COPY go.* ./
RUN go mod download
COPY . ./
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -gcflags=all="-l -B" -o /ssl-certify


FROM go AS dev
RUN apk add curl
USER haproxy
RUN curl -fLo install.sh https://raw.githubusercontent.com/cosmtrek/air/master/install.sh \
    && chmod +x install.sh && sh install.sh
CMD /app/bin/air
EXPOSE 9022


FROM gcr.io/distroless/static-debian11
WORKDIR /app
USER 99:99
COPY --chown=haproxy:haproxy --chmod=755 --from=builder /ssl-certify /app/ssl-certify
COPY --chown=haproxy:haproxy --chmod=755 --from=builder /app/config /app/config
ENTRYPOINT ["/app/ssl-certify"]
HEALTHCHECK --interval=12s --timeout=12s --start-period=30s \
    CMD ["/app/ssl-certify", "ping"]
CMD ["start"]
EXPOSE 9022
