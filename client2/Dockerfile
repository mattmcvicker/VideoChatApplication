# use `node` image as the base
FROM node:12
# set the current working directory to /app
WORKDIR /usr/src/app
# copy the package.json and package-lock.json files to the work directory
COPY package*.json ./
# run npm install to install all dependent packages
RUN npm install
# copy your JavaScript source files
COPY . .
# declare that your server will listen on port 80
EXPOSE 80 443
# set the entrypoint command
CMD ["node", "server.js"]