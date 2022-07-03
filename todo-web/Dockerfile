# syntax=docker/dockerfile:1

# Stage 1: Build
FROM node:18-alpine AS build-env
ENV NODE_ENV=production

WORKDIR /app

COPY package.json package-lock.json ./

RUN npm install --omit=dev

COPY . .

RUN npm run build

# Stage 2: Runtime
FROM nginx:stable-alpine

COPY --from=build-env /app/build /usr/share/nginx/html
