FROM alpine:edge
MAINTAINER Rosco Pecoltran <https://github.com/roscopecoltran>

# build: docker build -t scraper:alpine -f scraper-alpine.dockerfile --no-cache .
# run: docker run --rm --host -ti -p 3000:3000 -v `pwd`:/app scraper:alpine

ARG GOPATH=${GOPATH:-"/go"}
ARG APK_INTERACTIVE=${APK_INTERACTIVE:-"bash nano tree"}
ARG APK_RUNTIME=${APK_RUNTIME:-"go git openssl ca-certificates"}
ARG APK_BUILD=${APK_BUILD:-"gcc g++ musl-dev gfortran lapack-dev openssl-dev oniguruma-dev"}

ENV APP_BASENAME=${APP_BASENAME:-"scraper"} \
    PATH="${GOPATH}/bin:/app:$PATH" \
    GOPATH=${GOPATH:-"/go"}

RUN \
        apk add --no-cache ${APK_RUNTIME} && \
    \
        apk add --no-cache --virtual=.interactive-dependencies ${APK_INTERACTIVE} && \
    \
        apk add --no-cache --virtual=.build-dependencies ${APK_BUILD} && \
    \
        mkdir -p /data/cache
#    \
#      apk del --no-cache --virtual=.build-dependencies && \

COPY . /go/src/github.com/roscopecoltran/scraper
WORKDIR /go/src/github.com/roscopecoltran/scraper

RUN \
    go get github.com/Masterminds/glide && \
    go get github.com/mitchellh/gox && \
    \
    go get golang.org/x/net/... && \
    go get github.com/qor/session && \
    go get github.com/qor/action_bar && \
    go get github.com/qor/help && \
    go get github.com/qor/qor && \
    go get github.com/qor/admin && \
    go get github.com/qor/serializable_meta && \
    go get github.com/qor/worker && \
    go get github.com/qor/sorting && \
    go get github.com/qor/roles && \
    go get github.com/qor/publish && \
    go get github.com/qor/publish2 && \
    go get github.com/qor/oss/... && \
    go get github.com/jinzhu/gorm/... && \
    go get github.com/go-sql-driver/mysql && \
    go get github.com/roscopecoltran/admin && \
    go get gopkg.in/olahol/melody.v1 && \
    go get github.com/googollee/go-socket.io && \
    go get github.com/go-fsnotify/fsnotify && \
    go get gopkg.in/redis.v5  && \
    go get gopkg.in/redis.v3 && \
    go get golang.org/x/oauth2 && \
    go get github.com/victortrac/disks3cache && \
    go get github.com/tsak/concurrent-csv-writer && \
    go get github.com/trustmaster/goflow && \
    go get github.com/qor/media && \
    go get github.com/peterbourgon/diskv && \
    go get github.com/oleiade/reflections && \
    go get github.com/klaidliadon/go-redis-cache && \
    go get github.com/klaidliadon/go-memcached && \
    go get github.com/klaidliadon/go-couch-cache && \
    go get github.com/jeevatkm/go-model && \
    go get github.com/if1live/staticfilecache && \
    go get github.com/iancoleman/strcase && \
    go get github.com/ak1t0/flame && \
    go get github.com/alexflint/go-restructure && \
    go get github.com/archivers-space/warc && \
    go get github.com/birkelund/boltdbcache && \
    go get github.com/cabify/go-couchdb && \
    go get github.com/cnf/structhash && \
    go get github.com/go-mangos/mangos && \
    go get github.com/go-resty/resty && \
    go get github.com/gregjones/httpcache && \
    go get golang.org/x/crypto/... && \
    go get github.com/google/jsonapi && \
    go get github.com/adamzy/cedar-go && \
    go get github.com/google/gopacket && \
    go get github.com/jdkato/prose/... && \
    go get github.com/syndtr/goleveldb/... && \
    go get github.com/myntra/pipeline && \
    go get github.com/vmihailenco/msgpack && \
    go get github.com/alessio/shellescape && \
    go get github.com/jeffail/tunny && \
    go get github.com/plar/go-adaptive-radix-tree && \
    go get github.com/kamilsk/semaphore && \
    go get github.com/rjeczalik/interfaces/... && \
    go get github.com/kamildrazkiewicz/go-flow && \
    go get github.com/benmanns/goworker && \
    go get github.com/ibmendoza/msgq && \
    go get github.com/patrickmn/go-cache && \
    \
    glide install --strip-vendor

    # gox -verbose -os="linux" -arch="amd64" -output="/app/{{.Dir}}" ./cmd/scraper-server

VOLUME ["/data"]

EXPOSE 3000 4000

CMD ["/bin/bash"]
# CMD ["/app/scraper-server","/app/conf.d/providers.list.json"]