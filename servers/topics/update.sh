docker pull kjmasumo/feud-ms
docker rm -f feud-ms
docker run -d --network customNet --name feud-ms kjmasumo/feud-ms
docker rm -f mongocontainer
docker run -d --network customNet --name mongocontainer mongo
exit