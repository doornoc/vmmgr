FROM node:16.13.0 AS builder
WORKDIR /builder
RUN apt update
RUN apt install git
RUN git clone https://github.com/doornoc/vmmgr-web
WORKDIR /builder/vmmgr-web
RUN git checkout develop
RUN npm install
RUN npm run build-prod

FROM nginx:latest
WORKDIR /usr/share/nginx/html
COPY --from=builder /builder/vmmgr-web/build ./
