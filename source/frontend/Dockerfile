FROM node:20-alpine AS build
WORKDIR /app

ENV VITE_USE_MOCKS=false
ENV VITE_USE_MOCKS_MOVIES=false
ENV VITE_USE_MOCKS_REVIEWS=false
ENV VITE_USE_MOCKS_PLAYLISTS=false

COPY package*.json ./
RUN npm ci

COPY . .
RUN npm run build

FROM nginx:1.27-alpine

ENV NGINX_ENVSUBST_TEMPLATE_DIR=/etc/nginx/templates
ENV NGINX_ENVSUBST_OUTPUT_DIR=/etc/nginx/conf.d
ENV API_GATEWAY_URL=http://api-gateway:8080

COPY nginx.conf.template /etc/nginx/templates/default.conf.template
COPY --from=build /app/dist /usr/share/nginx/html
