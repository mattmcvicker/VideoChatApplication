# build go executable for linux
GOOS=linux go build

docker build -t kjmasumo/servers .
cd ../db
docker build -t kjmasumo/db .
cd ../gateway

# delete Go executable
go clean

