docker build -t kjmasumo/topics .
docker push kjmasumo/topics

ssh ec2-user@api.kenmasumoto.me < update.sh