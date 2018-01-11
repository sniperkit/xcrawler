####################################
# Build GOLANG based executables
####################################

### PLATFORM ########################################################################################################
ifeq (Darwin, $(findstring Darwin, $(shell uname -a)))
  PLATFORM := OSX
  OS := darwin
else
  PLATFORM := Linux
  OS := linux
endif

### ARCH ############################################################################################################
ARCH := amd64

### RUN #############################################################################################################
.PHONY: run
run:
	@go run main.go 

### BUILD ###########################################################################################################
.PHONY: build
build:
	@go build -o $(CURDIR)/dist/e3w-local main.go

### DIST ############################################################################################################
.PHONY: dist
dist:
	gox -verbose -os="darwin linux" -arch="amd64" -output="./dist/e3w-{{.OS}}" $(glide novendor)

### DEPS ############################################################################################################
.PHONY: deps
deps:
	@go get -v -u github.com/Masterminds/glide
	@go get -v -u github.com/mitchellh/gox

### DOCKER-COMPOSE ##################################################################################################
.PHONY: compose-dev
compose-dev:
	@docker-compose run --remove-orphans e3w-dev

.PHONY: compose
compose:
	@docker-compose up --remove-orphans e3w-dist