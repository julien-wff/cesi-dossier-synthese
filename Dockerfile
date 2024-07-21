FROM oven/bun:1-alpine AS web-builder

WORKDIR /app

COPY web/package.json web/bun.lockb ./
RUN bun install --frozen-lockfile && rm -rf /root/.bun

COPY web/ ./
ENV PUBLIC_API_ENDPOINT='/api'
RUN bun run build


FROM golang:1-alpine AS backend-builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
COPY --from=web-builder /app/build/ ./internal/web/build/
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/app ./cmd/server/main.go


FROM alpine:3.20 AS runner

USER 1000:1000

WORKDIR /app

# pdfcpu needs this
ENV XDG_CONFIG_HOME=/app/.config
RUN mkdir -p /app/.config

COPY --from=backend-builder /app/bin/app ./

EXPOSE 8080

CMD ["./app"]
