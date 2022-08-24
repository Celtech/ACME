FROM golang:1.16-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go get -u github.com/swaggo/swag/cmd/swag && go mod download
COPY . ./
RUN swag init --parseDependency --parseInternal --md docs --codeExampleFiles docs/samples && \
    CGO_ENABLED=0 go build -ldflags="-s -w -X main.Version=0.0.1" -gcflags=all="-l -B" -o /baker-acme

FROM golang:1.16-alpine AS base
WORKDIR /app
RUN addgroup --gid 99 --system haproxy && \
    adduser --disabled-password --home /var/lib/haproxy --ingroup haproxy --no-create-home --system --uid 99 haproxy && \
    mkdir /var/lib/haproxy && \
    chown haproxy:haproxy /var/lib/haproxy && \
    mkdir /data && \
    chown -R haproxy:haproxy /data && \
    chown -R haproxy:haproxy /app
EXPOSE 9022


FROM base AS dev
RUN apk add curl
USER haproxy
RUN curl -fLo install.sh https://raw.githubusercontent.com/cosmtrek/air/master/install.sh \
    && chmod +x install.sh && sh install.sh
CMD /app/bin/air


FROM base
COPY --chown=haproxy:haproxy --chmod=755 --from=builder /baker-acme /app/baker-acme
COPY --chown=haproxy:haproxy --chmod=755 --from=builder /app/config /app/config
COPY --chown=haproxy:haproxy --chmod=755 --from=builder /app/docs/swagger.json /app/docs/swagger.json
COPY --chown=haproxy:haproxy --chmod=755 --from=builder /app/docs/swagger.yaml /app/docs/swagger.yaml
RUN ln -s /app/baker-acme /usr/bin/baker-acme
USER haproxy
CMD ["/app/baker-acme", "start"]
