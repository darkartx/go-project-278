# Build
FROM golang:1.24-alpine AS build
RUN apk add --no-cache git
WORKDIR /build/code

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
  go mod download

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-X code.commitHash=$(git rev-parse HEAD)" -o /build/app ./cmd/url_shortener

# App
FROM alpine:3.22

ENV DEBUG=false
ENV PORT=8080
ENV ROLLBAR_SERVER_ROOT=https://github.com/darkartx/go-project-278
EXPOSE 8080

WORKDIR /app

COPY --from=build /build/app /app/bin/app

COPY --from=build build/code/db/migrations /app/db/migrations
COPY --from=build /go/bin/goose /usr/local/bin/goose

COPY bin/run.sh /app/bin/run.sh
RUN chmod +x /app/bin/run.sh

CMD ["/app/bin/run.sh"]
