#!/bin/sh
APP="OpenPipelineIO"

# assets 폴더의 모든 에셋을 빌드전에 assets_vfsdata.go 파일로 생성한다.
go run assets/asset_generate.go

build() {
    GOOS=$1
    GOARCH=$2
    OUTPUT_DIR="./bin/${GOOS}_${GOARCH}"

    mkdir -p ${OUTPUT_DIR}
    go build -ldflags "-X main.SHA1VER=`git rev-parse HEAD` -X main.BUILDTIME=`date -u +%Y-%m-%dT%H:%M:%S`" -o ${OUTPUT_DIR}/${APP} *.go

    # 압축
    cd ${OUTPUT_DIR}
    tar -zcvf ../${APP}_${GOOS}_${GOARCH}.tgz .
    cd -
}

# OS 및 아키텍처별 빌드
build windows amd64
build linux amd64
build linux riscv64
build freebsd amd64
build darwin amd64
build darwin arm64 # Apple Silicon 지원

# 디렉토리 삭제
rm -rf ./bin/windows_amd64
rm -rf ./bin/linux_amd64
rm -rf ./bin/linux_riscv64
rm -rf ./bin/freebsd_amd64
rm -rf ./bin/darwin_amd64
rm -rf ./bin/darwin_arm64
