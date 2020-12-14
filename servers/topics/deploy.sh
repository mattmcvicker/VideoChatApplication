docker build -t kjmasumo/feud-ms .
docker push kjmasumo/feud-ms

ssh ec2-user@kenmasumoto.me < update.sh