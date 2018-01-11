# Sniperkit-Scraper - Docker stack
[To do]

## Intro
[WIP]

### Features
-

### Goals
1.
2.
3.

### Quick Start
```bash
go get -v github.com/roscopecoltran/scraper
cd $GOPATH/src/github.com/roscopecoltran/scraper
go run *.go ./providers.dev.yaml
```

#### Make

##### (DEV) Scraper
```bash
make build
make run
```

##### (DIST) Scraper
```bash
make dist
```

#### Crane
```bash
go get -v -u github.com/michaelsauter/crane
```

##### (DIST) Scraper
```bash
crane up dist
```

##### (DEV) Scraper
```bash
crane up dev
```

#### Docker-Compose
MacOSX: 
```bash
brew install docker
brew install docker-compose
```

#### (DIST) Scraper + ETCD3 / E3CH 
Bootsrap:
```bash
docker-compose build --no-cache scraper
docker-compose up scraper
```

Examples:
```bash
open http://localhost:3000/bing?query=dlib (bing search endpoint)
open http://localhost:3000/admin (scraper admin)
```

#### (DEV) Scraper + ETCD3 / E3CH 
Bootsrap:
```bash
docker-compose build --no-cache scraper_dev
docker-compose up scraper_dev
```

#### (DEV) Scraper + ETCD3 / E3CH + ELK
Bootsrap:
```bash
docker-compose build --no-cache scraper_elk
docker-compose up scraper_elk
```

Examples:
```bash
open http://localhost:8086/ (e3ch)
open http://localhost:5601/ (kibana v5.x)
```

#### ETCD3 / E3CH 
Bootsrap:
```bash
docker-compose build --no-cache e3w_dev
docker-compose up e3w_dev
```

Examples:
```bash
open http://localhost:8086/ (e3ch)
```

go run *.go --debug --verbose ./providers.dev.json

ip="ifconfig en0 | grep inet | awk '$1=="inet" {print $2}'"
socat TCP-LISTEN:6000,reuseaddr,fork UNIX-CLIENT:\"$DISPLAY\"
eg. docker run -e DISPLAY=192.168.0.2:0 gns3/xeyes
https://stackoverflow.com/questions/37826094/xt-error-cant-open-display-if-using-default-display

socat TCP-LISTEN:6000,reuseaddr,fork UNIX-CLIENT:\"$DISPLAY\"
docker run -e DISPLAY=192.168.0.2:0 jess/geary

## RabbitMQ (brew/osx)
Management Plugin enabled by default at http://localhost:15672

Bash completion has been installed to:
  /usr/local/etc/bash_completion.d

To have launchd start rabbitmq now and restart at login:
  brew services start rabbitmq
Or, if you don't want/need a background service you can just run:
  rabbitmq-server

## NSQ
To have launchd start nsq now and restart at login:
  brew services start nsq
Or, if you don't want/need a background service you can just run:
  nsqd -data-path=/usr/local/var/nsq

