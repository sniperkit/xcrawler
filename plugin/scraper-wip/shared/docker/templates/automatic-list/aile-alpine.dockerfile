FROM frolvlad/alpine-glibc:alpine-3.6
MAINTAINER Rosco Pecoltran <https://github.com/roscopecoltran>

# build: docker build -t aile:alpine -f aile-alpine.dockerfile --no-cache .
# run: docker run --rm -ti -p 2312:2312 -v `pwd`:/app aile:alpine

# ref: https://repo.continuum.io/miniconda/
# - stable/default:
# ARG CONDA_VERSION=${CONDA_VERSION:-"4.0.5"}
# ARG CONDA_MD5_CHECKSUM=${CONDA_MD5_CHECKSUM:-"42dac45eee5e58f05f37399adda45e85"}
# - latest (oct 2017)
ARG CONDA_VERSION=${CONDA_VERSION:-"4.3.27.1"}
ARG CONDA_MD5_CHECKSUM=${CONDA_MD5_CHECKSUM:-"7efba9cbe774169e36695564f197becb"}
ARG CONDA_DIR=${CONDA_DIR:-"/opt/conda"}
# golang (optional)
ARG GOPATH=${GOPATH:-"/tmp/go"}

ENV APP_BASENAME=${APP_BASENAME:-"aile"}

ENV CONDA_DIR="/opt/conda" \
    PATH="$CONDA_DIR/bin:/app/cmd/${APP_BASENAME}/cli:/app/cmd/${APP_BASENAME}/server:${GOPATH}/bin:$PATH" \
    PKG_CONFIG_PATH="/usr/lib/pkgconfig/:/usr/local/lib/pkgconfig/" \
    CFLAGS=-I/usr/lib/python2.7/site-packages/numpy/core/include \
    PYTHONDONTWRITEBYTECODE=${PYTHONDONTWRITEBYTECODE:-"1"}

# echo "http://dl-4.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories && \
# Install conda
RUN \
        apk add --no-cache --virtual=.build-dependencies wget ca-certificates bash && \
    \
        mkdir -p "$CONDA_DIR" && \
        wget "http://repo.continuum.io/miniconda/Miniconda2-${CONDA_VERSION}-Linux-x86_64.sh" -O miniconda.sh && \
        echo "$CONDA_MD5_CHECKSUM  miniconda.sh" | md5sum -c && \
        bash miniconda.sh -f -b -p "$CONDA_DIR" && \
        echo "export PATH=$CONDA_DIR/bin:\$PATH" > /etc/profile.d/conda.sh && \
        rm miniconda.sh && \
    \
        conda update --all --yes && \
        conda config --set auto_update_conda False && \
        rm -r "$CONDA_DIR/pkgs" && \
    \
        apk del --purge .build-dependencies && \
    \
        mkdir -p "$CONDA_DIR/locks" && \
        chmod 777 "$CONDA_DIR/locks"

ARG APK_INTERACTIVE=${APK_INTERACTIVE:-"bash nano tree"}
ARG APK_RUNTIME=${APK_RUNTIME:-"libstdc++ cython git libx11 openssl ca-certificates"}
ARG APK_BUILD=${APK_BUILD:-"g++ gfortran musl-dev cython-dev libx11-dev gcc linux-headers lapack-dev python2-dev openssl-dev"}
ARG PIP_DEPS=${PIP_DEPS:-"numpy pandas scipy scikit-learn"}

COPY requirements.txt /app/requirements.txt

WORKDIR /app

RUN apk add --no-cache ${APK_RUNTIME} && \
    \
        apk add --no-cache --virtual=.interactive-dependencies ${APK_INTERACTIVE} && \
    \
        apk add --no-cache --virtual=.build-dependencies ${APK_BUILD} && \
    \
        ln -s locale.h /usr/include/xlocale.h && \
    \
        pip install --upgrade pip setuptools && \   
	    pip install --no-cache --no-cache-dir -r /app/requirements.txt && \
    \
        find /usr/lib/python2.*/ -name 'tests' -exec rm -r '{}' + && \
        rm /usr/include/xlocale.h && \
        rm -r /root/.cache && \
    \
    mkdir -p /data

    # apk del .build-dependencies && \

COPY . /app
VOLUME ["/data"]
EXPOSE 2312

RUN python setup.py develop

CMD ["/bin/bash"]

### SNIPPETS #########################################################################################################
# run (with x11): docker run --rm -it -e DISPLAY -v $(pwd):/app -v /tmp/.X11-unix:/tmp/.X11-unix:ro -v $XAUTHORITY:/root/.Xauthority --net=host aile:alpine
