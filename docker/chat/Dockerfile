FROM golang:alpine as build

WORKDIR /app
COPY ./ /app

RUN CGO_ENABLED=0 go build -o chatService services/chat/cmd/main.go


FROM alpine:latest

WORKDIR /app
COPY --from=build /app/chatService /app/
RUN chmod +x /app/chatService
ENTRYPOINT /app/chatService
