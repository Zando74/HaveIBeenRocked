FROM node:22-alpine as build-stage
WORKDIR /app
COPY frontend/package.json ./
COPY frontend/package-lock.json ./
RUN npm install --silent
COPY frontend/. ./
RUN npm run build

FROM nginx:stable-alpine
COPY --from=build-stage /app/dist /etc/nginx/build
COPY infra/nginx/default.conf /etc/nginx/conf.d/default.conf