docker stop $(docker ps -a -q)
docker rm $(docker ps -a -q)
docker volume prune -f
docker container prune -f
docker-compose up -d