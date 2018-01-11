docker run -it --rm -p 8080:8080 -v /docker/host/dir:/var/letmegrpc/protos -e PROTO_FILE file.proto -e SERVICE_ADDRESS localhost:8888 rudiscz/letmegrpc
docker run -it --rm -p 8080:8080 rudiscz/letmegrpc
https://github.com/rpliva/letmegrpc-docker