FROM golang:1.16-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go build -o /errandboi

FROM alpine:3.12

WORKDIR /app

COPY --from=builder /errandboi .

EXPOSE 3000

ENTRYPOINT [ "./erranboi" ]

CMD [ "serve" ]