FROM golang:1.22.2-alpine AS build

WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download && go mod verify
RUN go mod download

COPY ./ ./
RUN CGO_ENABLED=0 go build -a -o ./bin/message-service ./cmd/message/app.go


FROM alpine

WORKDIR /app

COPY --from=build /app/config ./config
COPY --from=build /app/bin ./

ENTRYPOINT [ "/app/message-service" ]