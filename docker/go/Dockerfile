FROM golang:1.14.2-alpine as build

WORKDIR /go/app

COPY . .

RUN set -ex \
    && go build -o build/crawler ./cmd/crawler/main.go


FROM alpine

WORKDIR /app

ENV APP_ENV=production

COPY --from=build /go/app/build/crawler .
COPY --from=build /go/app/configs/dbconf.yml ./configs/

RUN addgroup go \
    && adduser -D -G go go \
    && chown -R go:go /app/crawler \
    && chown -R go:go /app/configs/dbconf.yml

CMD ["./crawler"]