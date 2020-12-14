./build.sh

docker push kjmasumo/servers
docker push kjmasumo/db

ssh ec2-user@api.kenmasumoto.me < update.sh