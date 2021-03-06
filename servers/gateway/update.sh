docker pull kjmasumo/servers
docker pull kjmasumo/db

# environment variables
export SESSIONKEY="testsessionkey"
export REDISADDR=rServe:6379
export MYSQL_ROOT_PASSWORD=sqlpassword
export DSN=root:sqlpassword@tcp\(mysql:3306\)/mysql

docker rm -f api-server
docker rm -f rServe
docker rm -f mysql

# create mysql container
docker run -d \
--network customNet \
--name mysql \
-e MYSQL_ROOT_PASSWORD=sqlpassword \
-e MYSQL_DATABASE=mysql \
kjmasumo/db

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
-v /etc/letsencrypt:/etc/letsencrypt:ro \
kjmasumo/servers