FROM golang:alpine as build

WORKDIR /app
COPY ./ /app

RUN CGO_ENABLED=0 go build -o authService services/auth/cmd/main.go


FROM alpine:latest

WORKDIR /app
COPY --from=build /app/authService /app/
RUN chmod +x /app/authService
ENTRYPOINT /app/authService
