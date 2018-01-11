
## App - Default Targets
SWAGGER_UI_VERSION=3.3.2

run:
	# @go run *.go $(CURDIR)/shared/conf.d/providers.json
	@go run *.go $(CURDIR)/shared/conf.d/providers.yaml

build:
	@go build -o $(CURDIR)/dist/scraper-local main.go
	@echo "$ ./dist/scraper-local ./shared/conf.d/providers.yaml"
	@echo ""
	@./dist/scraper-local ./shared/conf.d/providers.yaml

dist:
	gox -verbose -os="darwin linux" -arch="amd64" -output="./dist/scraper-{{.OS}}" $(glide novendor)
	# gox -verbose -os="darwin linux" -arch="amd64" -output="./dist/scraper-{{.OS}}" $(glide novendor)

deps:
	@go get -v -u github.com/Masterminds/glide
	@go get -v -u github.com/mitchellh/gox
	@go get -v -u github.com/moovweb/rubex

compose:
	@docker-compose up --remove-orphans scraper

swagger-ui:	
	curl -L -o $(CURDIR)/swagger-ui-${SWAGGER_UI_VERSION}.tar.gz https://github.com/swagger-api/swagger-ui/archive/v$(SWAGGER_UI_VERSION).tar.gz
	tar zxf $(CURDIR)/swagger-ui-$(SWAGGER_UI_VERSION).tar.gz
	mv $(CURDIR)/swagger-ui-$(SWAGGER_UI_VERSION) $(CURDIR)/swaggerui
	rm -f $(CURDIR)/swagger-ui-$(SWAGGER_UI_VERSION).tar.gz

## App - Media Processing
## ref. https://github.com/mohanson/FaceDetectionServer
SrcDir=/src
SeetaCommitID=0f73c0964cf229d16fe584db14c08c61b1d84105
SeetaFDSrcDir=$(SrcDir)/SeetaFaceEngine/FaceDetection
CMake3=cmake3

seeta/clone:
	mkdir -p $(SrcDir)
	cd $(SrcDir) && git clone https://github.com/seetaface/SeetaFaceEngine.git
	cd $(SrcDir)/SeetaFaceEngine && git checkout $(SeetaCommitID)
.PHONY: seeta/clone

seeta/build:
	rm -rf $(SeetaFDSrcDir)/build; mkdir $(SeetaFDSrcDir)/build
	cd $(SeetaFDSrcDir)/build; $(CMake3) ..; make -j${nproc}
	cp $(SeetaFDSrcDir)/build/libseeta_facedet_lib.so /lib64/libseeta_facedet_lib.so
.PHONY: seeta/build

seeta: seeta/clone seeta/build
.PHONY: seeta

faced:
	cd libfaced && g++ -std=c++11 faced.cpp -fPIC -shared -o libfaced.so `pkg-config opencv --cflags --libs` \
		-I$(SeetaFDSrcDir)/include/ \
		-L$(SeetaFDSrcDir)/build \
		-lseeta_facedet_lib -ljsoncpp
	cd libfaced && cp libfaced.so /lib64/libfaced.so
	cd libfaced && g++ -std=c++11 faced_cmd.cpp -o faced \
		-I$(SeetaFDSrcDir)/include \
		-L$(SeetaFDSrcDir)/build \
		-L. -lfaced
	cd libfaced && ./faced ../face.jpg
.PHONY: faced

goserver:
	go build server.go
.PHONY: goserver

clean:
	rm -f /lib64/libfaced_facedet_lib.so
	rm -f /lib64/libfaced.so
	rm -f libfaced/faced
	rm -f libfaced/libfaced.so
.PHONY: clean

