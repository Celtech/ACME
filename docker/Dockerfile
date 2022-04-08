FROM golang:1.16-alpine AS base
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -gcflags=all="-l -B" -o /challenge-server


FROM scratch
COPY --chmod=755 --from=base /challenge-server /challenge-server
EXPOSE 9022
CMD [ "/challenge-server" ]
