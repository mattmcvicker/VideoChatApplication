# build go executable for linux
# possible run GOOS=linux go build -o gateway
GOOS=linux go build

# TODO: decide on docker image 
docker build -t kjmasumo/servers .
cd ../db
docker build -t kjmasumo/db .
cd ../gateway

# delete Go executable
go clean

