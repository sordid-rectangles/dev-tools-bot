FROM golang:1.17-alpine AS build

RUN mkdir /build

ADD . /build/

WORKDIR /build

RUN go mod download && go mod verify

RUN CGO_ENABLED=0 GOOS=linux go build -a -o bot .



FROM alpine:latest 

RUN mkdir /app
COPY --from=build /build/bot ./app

ENTRYPOINT [ "./app/bot" ]




