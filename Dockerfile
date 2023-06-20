FROM golang:1.20 AS builder

WORKDIR /app

COPY go.mod .
RUN go mod download

COPY . .

RUN go build -o auth ./src/.

FROM alpine:latest

RUN apk add libc6-compat

WORKDIR /app
COPY --from=builder /app/auth /app/auth
COPY app.env app.env

EXPOSE 9993

CMD [ "/app/auth" ]
