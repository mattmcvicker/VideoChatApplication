docker pull kjmasumo/topics
docker rm -f mongocontainer
docker run -d --network customNet --name mongocontainer mongo
docker rm -f topics
docker run -d --network customNet --name topics kjmasumo/topics
exit