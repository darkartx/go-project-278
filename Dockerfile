# Fontend build
FROM node:24-alpine AS frontend_builder
WORKDIR /build/frontend

RUN --mount=type=cache,target=/root/.npm \
  npm i @hexlet/project-url-shortener-frontend && \
  npm ci --prefer-offline --no-audit

# Backend build
FROM golang:1.24-alpine AS backend_builded
RUN apk add --no-cache git
WORKDIR /build/code

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
  go mod download

RUN go install github.com/pressly/goose/v3/cmd/goose@v3.26

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-X code.commitHash=$(git rev-parse HEAD)" -o /build/app .

# App
FROM alpine:3.22

ENV DEBUG=false
ENV ROLLBAR_SERVER_ROOT=https://github.com/darkartx/go-project-278
EXPOSE 80

RUN apk add --no-cache ca-certificates tzdata bash caddy

WORKDIR /app

COPY --from=frontend_builder \
  /build/frontend/node_modules/@hexlet/project-url-shortener-frontend/dist \
  /app/public
COPY --from=backend_builded /build/app /app/bin/app
COPY --from=backend_builded /build/code/db/migrations /app/db/migrations
COPY --from=backend_builded /go/bin/goose /usr/local/bin/goose
COPY Caddyfile /etc/caddy/Caddyfile

COPY bin/run.sh /app/bin/run.sh
RUN chmod +x /app/bin/run.sh

CMD ["/app/bin/run.sh"]
