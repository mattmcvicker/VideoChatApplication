docker build -t kjmasumo/client .
docker push kjmasumo/client

ssh ec2-user@kenmasumoto.me < update.sh