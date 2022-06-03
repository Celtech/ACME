FROM golang:1.16-alpine AS base
WORKDIR /app

COPY go.mod go.sum ./
RUN go get -u github.com/swaggo/swag/cmd/swag && go mod download

COPY . ./

RUN swag init --md docs --codeExampleFiles docs/samples && \
    CGO_ENABLED=0 go build -ldflags="-s -w" -gcflags=all="-l -B" -o /baker-acme



FROM openapitools/openapi-generator-cli:latest-release AS openapi
COPY --chown=haproxy:haproxy --chmod=755 --from=base /app/docs/swagger.json /app/docs/swagger.json
COPY --chown=haproxy:haproxy --chmod=755 --from=base /app/docs/swagger.yaml /app/docs/swagger.yaml
RUN ./usr/local/bin/docker-entrypoint.sh generate -i /app/docs/swagger.yaml -o /app/docs/openapi -g openapi-yaml --minimal-update && \
    ./usr/local/bin/docker-entrypoint.sh generate -s -i /app/docs/swagger.json -o /app/docs/openapi/openapi -g openapi --minimal-update



FROM golang:1.16-alpine
WORKDIR /app
RUN addgroup --gid 99 --system haproxy && \
    adduser --disabled-password --home /var/lib/haproxy --ingroup haproxy --no-create-home --system --uid 99 haproxy && \
    mkdir /var/lib/haproxy && \
    chown haproxy:haproxy /var/lib/haproxy
COPY --chown=haproxy:haproxy --chmod=755 --from=base /baker-acme /app/baker-acme
COPY --chown=haproxy:haproxy --chmod=755 --from=base /app/config /app/config
COPY --chown=haproxy:haproxy --chmod=755 --from=openapi /app/docs/openapi/openapi/openapi.json /app/docs/swagger.json
COPY --chown=haproxy:haproxy --chmod=755 --from=openapi /app/docs/openapi/openapi/openapi.yaml /app/docs/swagger.yaml
USER haproxy
EXPOSE 9022
CMD [ "/app/baker-acme" ]
