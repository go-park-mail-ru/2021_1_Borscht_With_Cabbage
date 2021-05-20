FROM golang:alpine as build

WORKDIR /app
COPY ./ /app

RUN CGO_ENABLED=0 go build -o mainService app/main.go


FROM alpine:latest

WORKDIR /app
COPY --from=build /app/mainService /app/
RUN chmod +x /app/mainService
ENTRYPOINT /app/mainService
