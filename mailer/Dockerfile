FROM golang:1.23.5-alpine AS build

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main ./cmd/mailer/main.go

FROM alpine:latest AS runner

WORKDIR /mailer

COPY --from=build /build/main ./main

CMD ["./main"]