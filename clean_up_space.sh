#docker stop $(docker ps -a -q)
#docker rm $(docker ps -a -q)
docker stop news-scraper-workers-daily-go
docker stop ews-scraper-workers-historic-go
docker rm news-scraper-workers-daily-go
docker rm ews-scraper-workers-historic-go
docker volume prune -f
docker container prune -f
docker-compose up -d