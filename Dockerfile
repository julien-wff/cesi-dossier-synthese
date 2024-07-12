FROM node:18.14-alpine AS web

RUN apk add --no-cache openjdk8-jre && npm install -g pnpm

WORKDIR /app/web

COPY web/package.json web/pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile

COPY . ./
RUN pnpm build

CMD ["node", "build"]
