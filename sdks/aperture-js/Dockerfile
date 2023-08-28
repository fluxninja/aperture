FROM node:16.19.1

# Create app directory
WORKDIR /usr/src/app

# Install app dependencies
# A wildcard is used to ensure both package.json AND package-lock.json are copied
# where available (npm@5+)
COPY package*.json tsconfig.json ./

RUN npm ci
# If you are building your code for production
# RUN npm ci --only=production

# Bundle app source
COPY . .

RUN npm run pre-build && npm run build && npm run post-build

RUN npm link

RUN cd example && npm ci && npm link @fluxninja/aperture-js && npm run build

HEALTHCHECK --interval=5s --timeout=60s --retries=3 --start-period=5s \
   CMD wget --no-verbose --tries=1 --spider 127.0.0.1:8080/health || exit 1

CMD [ "node", "./example/dist/example.js" ]
