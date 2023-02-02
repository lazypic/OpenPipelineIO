#!/bin/sh
APP="OpenPipelineIO"

# assets 폴더의 모든 에셋을 빌드전에 assets_vfsdata.go 파일로 생성한다.
go run assets/asset_generate.go

# OS별 기본빌드
GOOS=windows GOARCH=amd64 go build -ldflags "-X main.SHA1VER=`git rev-parse HEAD` -X main.BUILDTIME=`date -u +%Y-%m-%dT%H:%M:%S`" -o ./bin/windows/${APP} *.go
GOOS=linux GOARCH=amd64 go build -ldflags "-X main.SHA1VER=`git rev-parse HEAD` -X main.BUILDTIME=`date -u +%Y-%m-%dT%H:%M:%S`" -o ./bin/linux/${APP} *.go
GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.SHA1VER=`git rev-parse HEAD` -X main.BUILDTIME=`date -u +%Y-%m-%dT%H:%M:%S`" -o ./bin/darwin/${APP} *.go

# Github Release에 업로드 하기위해 압축
cd ./bin/windows/ && tar -zcvf ../${APP}_windows_x86-64.tgz . && cd -
cd ./bin/linux/ && tar -zcvf ../${APP}_linux_x86-64.tgz . && cd -
cd ./bin/darwin/ && tar -zcvf ../${APP}_darwin_x86-64.tgz . && cd -

# 삭제
rm -rf ./bin/windows
rm -rf ./bin/linux
rm -rf ./bin/darwin