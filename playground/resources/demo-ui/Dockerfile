FROM node:18-alpine

WORKDIR /src

COPY package.json package.json
COPY package-lock.json package-lock.json
RUN npm ci

COPY . /src

CMD npm start
