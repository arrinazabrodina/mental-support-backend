FROM node:21-alpine

LABEL authors="arinazabrodina"

WORKDIR /usr/src

# Expose port for service
EXPOSE 3000

# Install and configure `serve`.
RUN npm install -g serve

# Copy source code to image
COPY . .

# Install dependencies
RUN npm install

# Build app and start server from script
CMD ["/usr/src/run"]

#WORKDIR /app
#
#COPY . .
#
#RUN npm ci
#
#RUN npm run build
#
#ENV NODE_ENV production
#
#EXPOSE 3000
#
#CMD [ "npx", "serve", "build" ]