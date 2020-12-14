# build mongodb container
docker rm -f mongocontainer
docker run -d --network customNet --name mongocontainer mongo

# build nodejs container
docker build -t kjmasumo/feud-ms .
docker rm -f feud-ms
docker run -d --network customNet --name feud-ms kjmasumo/feud-ms