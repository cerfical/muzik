FROM golang:alpine AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY cmd/api/main.go main.go
COPY internal internal
RUN go build -tags viper_bind_struct -v -o bin/app main.go

FROM alpine:latest

COPY --from=build /app/bin/ /usr/local/bin/

ENV MUZIK_SERVER_ADDR=:80
EXPOSE 80

ENTRYPOINT ["app"]
