# node is an image that already has nodejs installed on it
FROM node

# set the current working directory in container
WORKDIR /app

COPY package.json package-lock.json ./

RUN npm install

COPY . .

EXPOSE 80
ENTRYPOINT ["node", "index.js"]