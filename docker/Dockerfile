FROM golang:1.16-alpine AS base
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -gcflags=all="-l -B" -o /baker-acme


FROM golang:1.16-alpine
COPY --chmod=755 --from=base /baker-acme /baker-acme
EXPOSE 9022
CMD [ "/baker-acme" ]
