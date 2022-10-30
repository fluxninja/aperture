FROM node:16

# Create app directory
WORKDIR /usr/src/app

# Install app dependencies
# A wildcard is used to ensure both package.json AND package-lock.json are copied
# where available (npm@5+)
COPY package*.json ./

RUN npm install
# If you are building your code for production
# RUN npm ci --only=production

# Bundle app source
COPY . .

HEALTHCHECK --interval=5s --timeout=60s --retries=3 --start-period=5s \
    CMD wget --no-verbose --tries=1 --spider 127.0.0.1:8080/health || exit 1

CMD [ "node", "./example/example.js" ]
