docker stop webhook
docker build . -t webhook
docker run -d -p 8080:8080 --name webhook webhook