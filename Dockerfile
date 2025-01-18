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


FROM alpine:3.21 AS runtime

USER 1000:1000

WORKDIR /app

# pdfcpu needs this
ENV XDG_CONFIG_HOME=/app/.config
RUN mkdir -p /app/.config

COPY --from=backend-builder /app/bin/app ./

ARG BUILD_DATE
ARG VCS_REF
ARG VCS_URL
ARG VERSION

LABEL org.opencontainers.image.version=$VERSION \
      org.opencontainers.image.title="CESI Dossier de Synthèse" \
      org.opencontainers.image.description="Web interface to visualize grades from CESI Engineering School" \
      org.opencontainers.image.authors="Julien W <cefadrom1@gmail.com>" \
      org.opencontainers.image.url=$VCS_URL \
      org.opencontainers.image.source=$VCS_URL \
      org.opencontainers.image.revision=$VCS_REF \
      org.opencontainers.image.created=$BUILD_DATE

EXPOSE 8080

CMD ["./app"]
