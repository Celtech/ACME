FROM golang:1.16-alpine AS base

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN go build -o /challenge-server


FROM alpine:3.15.4
COPY --chmod=755 --from=base /challenge-server /challenge-server
EXPOSE 9022
CMD [ "/challenge-server" ]
