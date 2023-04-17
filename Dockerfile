FROM node:18.14-alpine AS web

RUN apk add --no-cache openjdk8-jre && npm install -g pnpm

WORKDIR /app

COPY package.json pnpm-lock.yaml client/package.json server/package.json ./
RUN pnpm install --frozen-lockfile

COPY . ./
RUN pnpm build

CMD ["node", "build"]
