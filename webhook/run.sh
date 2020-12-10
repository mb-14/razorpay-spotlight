docker stop webhook
docker container rm webhook
docker build . -t webhook
docker run -d --net=host --name webhook webhook