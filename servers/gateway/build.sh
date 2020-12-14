# build go executable for linux
# possible run GOOS=linux go build -o gateway
GOOS=linux go build

# TODO: decide on docker image 
docker build -t kjmasumo/servers .
cd ../db
docker build -t kjmasumo/db .
cd ../gateway

# delete Go executable
go clean

# environment variables
export SESSIONKEY="testsessionkey"
export REDISADDR=rServe:6379
export MYSQL_ROOT_PASSWORD=sqlpassword
export DSN=root:sqlpassword@tcp\(mysql:3306\)/mysql
export FEUDADDR=fued-ms

docker rm -f api-server
docker rm -f rServe
docker rm -f mysql
docker network rm customNet

docker network create customNet

# create mysql container
docker run -d \
--network customNet \
--name mysql \
-e MYSQL_ROOT_PASSWORD=sqlpassword \
-e MYSQL_DATABASE=mysql \
kjmasumo/db # TODO: decide on docker image 

# create redis container
docker run  -d --rm --network customNet --name rServe redis

# run docker container
docker run -d -p 443:443 \
--name api-server \
--network customNet \
-e ADDR=:443 \
-e TLSCERT=/etc/letsencrypt/live/api.kenmasumoto.me/fullchain.pem \
-e TLSKEY=/etc/letsencrypt/live/api.kenmasumoto.me/privkey.pem \
-e REDISADDR=$REDISADDR \
-e SESSIONKEY=$SESSIONKEY \
-e DSN=$DSN \
-e FEUDADDR=$FEUDADDR \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
kjmasumo/servers # TODO: decide on docker image 
