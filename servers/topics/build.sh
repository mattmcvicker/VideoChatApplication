# build mongodb container
docker rm -f mongocontainer
docker run -d --network customNet --name mongocontainer mongo

# build nodejs container
docker build -t kjmasumo/topics .
docker rm -f topics
docker run -d --network customNet --name topics kjmasumo/topics