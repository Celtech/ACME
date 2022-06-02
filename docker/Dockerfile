FROM golang:1.16-alpine AS base
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -gcflags=all="-l -B" -o /baker-acme


FROM golang:1.16-alpine
RUN addgroup --gid 99 --system haproxy && \
    adduser --disabled-password --home /var/lib/haproxy --ingroup haproxy --no-create-home --system --uid 99 haproxy && \
    mkdir /var/lib/haproxy && \
    chown haproxy:haproxy /var/lib/haproxy
COPY --chown=haproxy:haproxy --chmod=755 --from=base /baker-acme /baker-acme
USER haproxy
EXPOSE 9022
CMD [ "/baker-acme" ]
