# TODO: decide on docker image
docker build -t {DOCKER_ACCOUNT}/db .

# run docker container
docker run -d -p 3306:3306 \
--network customNet \
--name mysql \
-e MYSQL_ROOT_PASSWORD=sqlpassword \
-e MYSQL_DATABASE=mysql \
{DOCKER_ACCOUNT}/db # TODO: decide on docker image