FROM golang:1.13 as build

WORKDIR /app

COPY . .
RUN CGO_ENABLED=0 go build main.go

FROM alpine

RUN apk update \
        && apk upgrade \
        && apk add --no-cache \
        ca-certificates \
        && update-ca-certificates 2>/dev/null || true

WORKDIR /app
COPY --from=build /app/main /app/main

RUN chmod +x /app/main

CMD [ "./main" ]