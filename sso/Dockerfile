FROM golang:1.23.5-alpine AS build

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main ./cmd/sso/main.go

FROM alpine:latest AS runner

WORKDIR /auth

COPY --from=build /build/main ./main
COPY --from=build /build/config/config_local.yaml ./config_local.yaml

CMD ["./main", "--config=./config_local.yaml"]