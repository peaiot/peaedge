BUILD_ORG   := toughstruct
BUILD_VERSION   := latest
BUILD_TIME      := $(shell date "+%F %T")
BUILD_NAME      := peabss
RELEASE_VERSION := v1.0.1
SOURCE          := main.go
RELEASE_DIR     := ./release
COMMIT_SHA1     := $(shell git show -s --format=%H )
COMMIT_DATE     := $(shell git show -s --format=%cD )
COMMIT_USER     := $(shell git show -s --format=%ce )
COMMIT_SUBJECT     := $(shell git show -s --format=%s )

clean:
	rm -f peabss

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
	buildpre
	go generate
	CGO_ENABLED=0 go build -a -ldflags \
	'\
	-X "main.BuildVersion=${BUILD_VERSION}"\
	-X "main.ReleaseVersion=${RELEASE_VERSION}"\
	-X "main.BuildTime=${BUILD_TIME}"\
	-X "main.BuildName=${BUILD_NAME}"\
	-X "main.CommitID=${COMMIT_SHA1}"\
	-X "main.CommitDate=${COMMIT_DATE}"\
	-X "main.CommitUser=${COMMIT_USER}"\
	-X "main.CommitSubject=${COMMIT_SUBJECT}"\
	-s -w -extldflags "-static"\
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
    -o ${RELEASE_DIR}/${BUILD_NAME} ${SOURCE}


ci:
	@read -p "type commit message: " cimsg; \
	git ci -am "$(shell date "+%F %T") $${cimsg}"


syncwjt:
	@read -p "提示:同步操作尽量在完成一个完整功能特性后进行，请输入提交描述 (wjt):  " cimsg; \
	git commit -am "$(shell date "+%F %T") : $${cimsg}" || echo "no commit"
	# 切换主分支并更新
	git checkout develop
	git pull origin develop
	# 切换开发分支变基合并提交
	git checkout wjt
	git rebase -i develop
	# 切换回主分支并合并开发者分支，推送主分支到远程，方便其他开发者合并
	git checkout develop
	git merge --no-ff wjt
	git push origin develop
	# 切换回自己的开发分支继续工作
	git checkout wjt


.PHONY: clean build


