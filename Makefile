BUILD_ORG   := toughstruct
BUILD_VERSION   := latest
BUILD_TIME      := $(shell date "+%F %T")
BUILD_NAME      := peaedge
RELEASE_VERSION := v1.0.1
SOURCE          := main.go
RELEASE_DIR     := ./release
COMMIT_SHA1     := $(shell git show -s --format=%H )
COMMIT_DATE     := $(shell git show -s --format=%cD )
COMMIT_USER     := $(shell git show -s --format=%ce )
COMMIT_SUBJECT     := $(shell git show -s --format=%s )

clean:
	rm -f peaedge

buildpre:
	echo "${BUILD_VERSION} ${RELEASE_VERSION} ${BUILD_TIME}" > assets/buildver.txt
	echo "BuildVersion=${BUILD_VERSION}" > assets/build.txt
	echo "ReleaseVersion=${RELEASE_VERSION}" >> assets/build.txt
	echo "BuildTime=${BUILD_TIME}" >> assets/build.txt
	echo "BuildName=${BUILD_NAME}" >> assets/build.txt
	echo "CommitID=${COMMIT_SHA1}" >> assets/build.txt
	echo "CommitDate=${COMMIT_DATE}" >> assets/build.txt
	echo "CommitUser=${COMMIT_USER}" >> assets/build.txt
	echo "CommitSubject=${COMMIT_SUBJECT}" >> assets/build.txt

build:
	go generate
	make buildpre
	CGO_ENABLED=1 go build --tags "libsqlite3 darwin" -a -ldflags \
	'\
	-X "main.BuildVersion=${BUILD_VERSION}"\
	-X "main.ReleaseVersion=${RELEASE_VERSION}"\
	-X "main.BuildTime=${BUILD_TIME}"\
	-X "main.BuildName=${BUILD_NAME}"\
	-X "main.CommitID=${COMMIT_SHA1}"\
	-X "main.CommitDate=${COMMIT_DATE}"\
	-X "main.CommitUser=${COMMIT_USER}"\
	-X "main.CommitSubject=${COMMIT_SUBJECT}"\
	' \
    -o ${BUILD_NAME} ${SOURCE}

build-linux:
	go generate
	make buildpre
	CC=x86_64-linux-musl-gcc CXX=x86_64-linux-musl-g++ \
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -ldflags \
	'\
	-X "main.BuildVersion=${BUILD_VERSION}"\
	-X "main.ReleaseVersion=${RELEASE_VERSION}"\
	-X "main.BuildTime=${BUILD_TIME}"\
	-X "main.BuildName=${BUILD_NAME}"\
	-X "main.CommitID=${COMMIT_SHA1}"\
	-X "main.CommitDate=${COMMIT_DATE}"\
	-X "main.CommitUser=${COMMIT_USER}"\
	-X "main.CommitSubject=${COMMIT_SUBJECT}"\
	-s -w -extldflags "-static" -linkmode "external" \
	' \
    -o ${RELEASE_DIR}/${BUILD_NAME}-x86_64 ${SOURCE}

build-linux-arm64:
	go generate
	make buildpre
	CC=aarch64-linux-musl-gcc CXX=aarch64-linux-musl-g++ \
	CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build -a -ldflags \
	'\
	-X "main.BuildVersion=${BUILD_VERSION}"\
	-X "main.ReleaseVersion=${RELEASE_VERSION}"\
	-X "main.BuildTime=${BUILD_TIME}"\
	-X "main.BuildName=${BUILD_NAME}"\
	-X "main.CommitID=${COMMIT_SHA1}"\
	-X "main.CommitDate=${COMMIT_DATE}"\
	-X "main.CommitUser=${COMMIT_USER}"\
	-X "main.CommitSubject=${COMMIT_SUBJECT}"\
	-s -w -extldflags "-static" -linkmode "external" \
	' \
    -o ${RELEASE_DIR}/${BUILD_NAME}-arm64 ${SOURCE}


ci:
	@read -p "type commit message: " cimsg; \
	git ci -am "$(shell date "+%F %T") $${cimsg}"

pub:
	make build-linux
	make build-linux-arm64
	upx ./release/peaedge-x86_64
	upx ./release/peaedge-arm64
	./peapub -put -m peaedge -i  ./release/peaedge-x86_64  -v latest
	./peapub -put -m peaedge-x86_64 -i  ./release/peaedge-x86_64  -v latest
	./peapub -put -m peaedge-arm64 -i  ./release/peaedge-arm64  -v latest

syncwjt:
	@read -p "??????:???????????????????????????????????????????????????????????????????????????????????? (wjt):  " cimsg; \
	git commit -am "$(shell date "+%F %T") : $${cimsg}" || echo "no commit"
	# ????????????????????????
	git checkout main
	git pull origin main
	# ????????????????????????????????????
	git checkout wjt
	git rebase -i main
	# ???????????????????????????????????????????????????????????????????????????????????????????????????
	git checkout main
	git merge --no-ff wjt
	git push origin main
	# ??????????????????????????????????????????
	git checkout wjt


.PHONY: clean build


