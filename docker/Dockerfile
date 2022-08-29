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
ARG BUILD_VERSION
ARG BUILD_DATE
COPY go.* ./
RUN go mod download
COPY . ./
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w -X github.com/Celtech/ACME/cmd.Version=${BUILD_VERSION} -X 'github.com/Celtech/ACME/cmd.Date=${BUILD_DATE}'" \
    -gcflags=all="-l -B" -o /ssl-certify


FROM go AS dev
RUN apk add curl
USER haproxy
RUN curl -fLo install.sh https://raw.githubusercontent.com/cosmtrek/air/master/install.sh \
    && chmod +x install.sh && sh install.sh
CMD /app/bin/air
EXPOSE 9022


FROM gcr.io/distroless/static-debian11
ARG BUILD_VERSION
ARG BUILD_DATE
ARG BUILD_VCS_REF
LABEL org.label-schema.schema-version="1.0"
LABEL org.label-schema.build-date=$BUILD_DATE
LABEL org.label-schema.name="ssl-certify/ssl-certify"
LABEL org.label-schema.description="Thin Let's Encrypt ACME client and challenge server written in go."
LABEL org.label-schema.url="https://sslcertify.com/"
LABEL org.label-schema.vcs-url="https://github.com/Celtech/ssl-certify"
LABEL org.label-schema.vcs-ref=$BUILD_VCS_REF
LABEL org.label-schema.vendor="SSL Certify"
LABEL org.label-schema.version=$BUILD_VERSION
LABEL org.label-schema.docker.cmd="docker run -p 9022:9022 -d ssl-certify/ssl-certify"
WORKDIR /app
USER 99:99
COPY --chown=haproxy:haproxy --chmod=755 --from=builder /ssl-certify /app/ssl-certify
COPY --chown=haproxy:haproxy --chmod=755 --from=builder /app/config /app/config
ENTRYPOINT ["/app/ssl-certify"]
HEALTHCHECK --interval=30s --timeout=30s --start-period=30s \
    CMD ["/app/ssl-certify", "ping"]
CMD ["start"]
EXPOSE 9022
