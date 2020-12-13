docker pull kjmasumo/client
docker rm -f client 
docker run -d -p 443:443 -p 80:80 -v /etc/letsencrypt:/etc/letsencrypt:ro --network customNetwork --name client kjmasumo/client
exit