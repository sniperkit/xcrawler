#!/bin/bash

git subtree add --prefix shared/docker/templates/alpine/py2/miniconda2 https://github.com/frol/docker-alpine-miniconda2 master --squash
git subtree add --prefix shared/docker/templates/alpine/py3/miniconda3 https://github.com/frol/docker-alpine-miniconda3 master --squash
git subtree add --prefix shared/docker/templates/alpine/cpp_c/glibc https://github.com/frol/docker-alpine-glibc master --squash
git subtree add --prefix shared/docker/templates/alpine/py3/flask-restplus-server https://github.com/frol/flask-restplus-server-example master --squash
git subtree add --prefix shared/docker/templates/alpine/java/oraclejdk8 https://github.com/frol/docker-alpine-oraclejdk8 master --squash
git subtree add --prefix shared/docker/templates/alpine/cpp_c/gcc https://github.com/frol/docker-alpine-gcc master --squash
git subtree add --prefix shared/docker/templates/alpine/py3/machinelearning https://github.com/frol/docker-alpine-python-machinelearning master --squash
git subtree add --prefix shared/docker/templates/alpine/base/fpc https://github.com/frol/docker-alpine-fpc master --squash
git subtree add --prefix shared/docker/templates/alpine/py3/base https://github.com/frol/docker-alpine-python3 master --squash
git subtree add --prefix shared/docker/templates/alpine/nim/base https://github.com/frol/docker-alpine-nim master --squash
git subtree add --prefix shared/docker/templates/alpine/rust/base https://github.com/frol/docker-alpine-rust master --squash
git subtree add --prefix shared/docker/templates/alpine/ruby/base https://github.com/frol/docker-alpine-ruby master --squash
git subtree add --prefix shared/docker/templates/alpine/py2/base https://github.com/frol/docker-alpine-python2 master --squash
git subtree add --prefix shared/docker/templates/alpine/java/openjdk7 https://github.com/frol/docker-alpine-openjdk7 master --squash
git subtree add --prefix shared/docker/templates/alpine/mono/base https://github.com/frol/docker-alpine-mono master --squash
git subtree add --prefix shared/docker/templates/alpine/scala/base https://github.com/frol/docker-alpine-scala master --squash
git subtree add --prefix shared/docker/templates/alpine/golang/base https://github.com/frol/docker-alpine-go master --squash
git subtree add --prefix shared/docker/templates/alpine/cpp_c/golang https://github.com/frol/docker-alpine-gxx master --squash
git subtree add --prefix shared/docker/templates/alpine/py3/xgboost https://github.com/petronetto/machine-learning-alpine master --squash
git subtree add --prefix shared/docker/templates/alpine/py3/tensorflow-jupyter https://github.com/tatsushid/docker-alpine-py3-tensorflow-jupyter master --squash
git subtree add --prefix shared/docker/templates/alpine/py3/tensorflow-keras https://github.com/petronetto/tensorflow-alpine master --squash
git subtree add --prefix shared/docker/templates/alpine/py3/tensorflow-keras-v2 https://github.com/smizy/docker-keras-tensorflow master --squash
git subtree add --prefix shared/docker/templates/alpine/py3/tensorflow https://github.com/feisan/alpine-python3-tensorflow master --squash
git subtree add --prefix shared/docker/templates/ai-kernels https://github.com/lablup/backend.ai-kernels master --squash