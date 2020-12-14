# build go executable for linux
# possible run GOOS=linux go build -o gateway
GOOS=linux go build

# TODO: decide on docker image 
docker build -t {DOCKER_ACCOUNT}/servers .

# delete Go executable
go clean

# environment variables
export SESSIONKEY="testsessionkey"
export REDISADDR=rServe:6379
export MYSQL_ROOT_PASSWORD=sqlpassword
export DSN="root:sqlpassword@tcp(mysql:3306)/mysql"

docker rm -f api-server
docker rm -f rServe
docker rm -f mysql

# create mysql container
docker run -d \
--network customNet \
--name mysql \
-e MYSQL_ROOT_PASSWORD=sqlpassword \
-e MYSQL_DATABASE=mysql \
{DOCKER_ACCOUNT}/db # TODO: decide on docker image 

# create redis container
docker run  -d --rm --network customNet --name rServe redis

# run docker container
docker run -d -p 80:80 \
--name api-server \
--network customNet \
-e REDISADDR=$REDISADDR -e \
SESSIONKEY=$SESSIONKEY \
-e DSN=$DSN \
{DOCKER_ACCOUNT}/servers # TODO: decide on docker image 
