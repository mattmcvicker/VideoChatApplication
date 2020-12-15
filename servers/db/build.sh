docker build -t kjmasumo/db .

# run docker container
docker run -d -p 3306:3306 \
--network customNet \
--name mysql \
-e MYSQL_ROOT_PASSWORD=sqlpassword \
-e MYSQL_DATABASE=mysql \
kjmasumo/db