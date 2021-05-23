FROM golang:alpine as build

WORKDIR /app
COPY ./ /app

RUN CGO_ENABLED=0 go build -o basketService services/basket/cmd/main.go


FROM alpine:latest

WORKDIR /app
COPY --from=build /app/basketService /app/
RUN chmod +x /app/basketService
ENTRYPOINT /app/basketService
