docker build . -t webhook
docker run -i -t -p 8080:8080 webhook